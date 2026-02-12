<template>
  <div class="app">
    <header>
      <h1>CI/CD 管道管理系统</h1>
      <nav>
        <ul>
          <li><button class="back-button" @click="goBack">← 回到上一页</button></li>
          <li><a href="#" @click.prevent="activeTab = 'projects'">项目管理</a></li>
          <li><a href="#" @click.prevent="activeTab = 'templates'">模板管理</a></li>
          <li><a href="#" @click.prevent="activeTab = 'guide'">使用指南</a></li>
        </ul>
      </nav>
    </header>

    <main>
      <!-- 项目管理 -->
      <section v-if="activeTab === 'projects'">
        <h2>项目管理</h2>
        
        <!-- 创建项目表单 -->
        <div class="form-section">
          <h3>创建新项目</h3>
          <div class="form-group">
            <label for="projectName">项目名称</label>
            <input type="text" id="projectName" v-model="newProject.name" placeholder="输入项目名称">
          </div>
          <div class="form-group">
            <label for="projectPath">项目路径</label>
            <input type="text" id="projectPath" v-model="newProject.path" placeholder="输入项目路径，例如：test/go-project">
          </div>
          <button class="primary" @click="createProject">创建项目</button>
        </div>

        <!-- 项目列表 -->
        <div class="project-list">
          <h3>项目列表</h3>
          <div v-if="loading">加载中...</div>
          <div v-else-if="projects.length === 0">暂无项目</div>
          <div v-else>
            <div v-for="project in projects" :key="project.id" class="project-item">
              <h3>{{ project.name }}</h3>
              <p>路径: {{ project.path }}</p>
              <p>ID: {{ project.id }}</p>
              <div class="project-actions">
                <button class="primary" @click="analyzeTechStack(project.id)">分析技术栈</button>
                <button class="primary" @click="showPlatformSelector('generate', project.id)">生成管道</button>
                <button class="primary" @click="showPlatformSelector('execute', project.id)">执行管道</button>
                <button class="secondary" @click="viewExecutions(project.id)">查看执行</button>
                <button class="secondary" @click="analyzeOptimization(project.id)">分析优化</button>
                <button class="danger" @click="deleteProject(project.id)">删除项目</button>
              </div>
            </div>

            <!-- 平台选择对话框 -->
            <div v-if="showPlatformDialog" class="platform-dialog">
              <div class="dialog-content">
                <h3>选择平台</h3>
                <p>请选择要使用的 CI/CD 平台：</p>
                <div class="platform-options">
                  <label>
                    <input type="radio" v-model="selectedPlatform" value="mock">
                    Mock (模拟平台)
                  </label>
                  <label>
                    <input type="radio" v-model="selectedPlatform" value="github_actions">
                    GitHub Actions
                  </label>
                </div>
                <div class="template-section" v-if="templates.length > 0">
                  <p>请选择要使用的模板：</p>
                  <select v-model="selectedTemplateId" class="template-select">
                    <option value="0">默认模板（根据技术栈自动选择）</option>
                    <optgroup label="内置模板">
                      <option v-for="template in templates.filter(t => t.is_builtin)" :key="template.id" :value="template.id">
                        {{ template.platform }} - {{ template.language || '通用' }}
                      </option>
                    </optgroup>
                    <optgroup label="自定义模板">
                      <option v-for="template in templates.filter(t => !t.is_builtin)" :key="template.id" :value="template.id">
                        {{ template.platform }} - {{ template.language || '通用' }}
                      </option>
                    </optgroup>
                  </select>
                </div>
                <div class="dialog-actions">
                  <button class="secondary" @click="showPlatformDialog = false">取消</button>
                  <button class="primary" @click="confirmPlatformSelection">确定</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 技术栈分析 -->
      <section v-else-if="activeTab === 'tech-stack'">
        <h2>技术栈分析</h2>
        <div v-if="techStackResult">
          <h3>分析结果</h3>
          <div class="tech-stack-details">
            <div class="tech-stack-item">
              <strong>编程语言:</strong>
              <span>{{ techStackResult.language || '未知' }}</span>
            </div>
            <div class="tech-stack-item">
              <strong>框架:</strong>
              <span>{{ techStackResult.framework || '无' }}</span>
            </div>
            <div class="tech-stack-item">
              <strong>构建工具:</strong>
              <span>{{ techStackResult.build_tool || '无' }}</span>
            </div>
            <div class="tech-stack-item">
              <strong>测试框架:</strong>
              <span>{{ techStackResult.test_framework || '无' }}</span>
            </div>
            <div class="tech-stack-item">
              <strong>依赖包:</strong>
              <div class="dependencies-list">
                <div v-if="techStackResult.dependencies" class="dependencies-grid">
                  <div v-for="(version, dep) in techStackResult.dependencies" :key="dep" class="dependency-item">
                    <span class="dependency-name">{{ dep }}</span>
                    <span class="dependency-version">{{ version }}</span>
                  </div>
                </div>
                <div v-else>
                  无依赖包
                </div>
              </div>
            </div>
            <div class="tech-stack-item">
              <strong>相关文件:</strong>
              <div class="files-list">
                <div v-if="techStackResult.files && techStackResult.files.length > 0">
                  <ul>
                    <li v-for="file in techStackResult.files" :key="file">{{ file }}</li>
                  </ul>
                </div>
                <div v-else>
                  无相关文件
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-else>
          <p>请从项目管理页面选择项目进行技术栈分析</p>
        </div>
      </section>

      <!-- 管道配置 -->
      <section v-else-if="activeTab === 'pipeline'">
        <h2>管道配置</h2>
        <div v-if="pipelineResult">
          <h3>生成结果</h3>
          <pre>{{ typeof pipelineResult === 'string' ? pipelineResult : pipelineResult.content }}</pre>
        </div>
        <div v-else>
          <p>请从项目管理页面选择项目生成管道配置</p>
        </div>
      </section>

      <!-- 执行管理 -->
      <section v-else-if="activeTab === 'execution'">
        <h2>执行管理</h2>
        <div v-if="executions.length > 0">
          <h3>执行历史</h3>
          <div class="table-container">
            <table class="table">
              <thead>
                <tr>
                  <th>执行ID</th>
                  <th>项目ID</th>
                  <th>状态</th>
                  <th>平台</th>
                  <th>开始时间</th>
                  <th>结束时间</th>
                  <th>持续时间(秒)</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="execution in executions" :key="execution.id">
                  <td>{{ execution.id }}</td>
                  <td>{{ execution.project_id }}</td>
                  <td>{{ execution.status }}</td>
                  <td>{{ execution.platform }}</td>
                  <td>{{ execution.start_time }}</td>
                  <td>{{ execution.end_time }}</td>
                  <td>{{ execution.duration }}</td>
                  <td>
                    <button class="secondary" @click="getExecutionMetrics(execution.id)">查看指标</button>
                    <button class="secondary" @click="viewExecutionDetails(execution.id)">查看详情</button>
                    <button class="danger" @click="stopExecution(execution.id)">停止执行</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div v-else>
          <p>请从项目管理页面选择项目执行管道</p>
        </div>
      </section>

      <!-- 指标分析 -->
      <section v-else-if="activeTab === 'metrics'">
        <h2>指标分析</h2>
        <div v-if="projectMetrics">
          <h3>项目指标</h3>
          <div class="metrics-grid">
            <div class="metric-card">
              <h4>平均执行时间</h4>
              <div class="value">{{ projectMetrics.average_duration || 0 }}s</div>
            </div>
            <div class="metric-card">
              <h4>成功率</h4>
              <div class="value">{{ projectMetrics.success_rate || 0 }}%</div>
            </div>
            <div class="metric-card">
              <h4>失败率</h4>
              <div class="value">{{ projectMetrics.failure_rate || 0 }}%</div>
            </div>
          </div>
        </div>
        <div v-else>
          <p>请从项目管理页面选择项目查看指标</p>
        </div>
      </section>

      <!-- 优化建议 -->
      <section v-else-if="activeTab === 'optimization'">
        <h2>优化建议</h2>
        <div v-if="optimizationResult">
          <h3>分析结果</h3>
          <div v-if="optimizationResult.metrics">
            <h4>执行指标</h4>
            <pre>{{ JSON.stringify(optimizationResult.metrics, null, 2) }}</pre>
          </div>
          <div v-if="optimizationResult.suggestions">
            <h4>优化建议</h4>
            <ul>
              <li v-for="(suggestion, index) in optimizationResult.suggestions" :key="index">
                {{ suggestion.description }}: {{ suggestion.suggestion }}
              </li>
            </ul>
          </div>
        </div>
        <div v-else>
          <p>请从项目管理页面选择项目分析优化建议</p>
        </div>
      </section>

      <!-- 执行详情 -->
      <section v-else-if="activeTab === 'execution-details'">
        <h2>执行详情</h2>
        <div v-if="executionDetails">
          <div class="execution-info">
            <div class="info-item">
              <strong>执行ID:</strong>
              <span>{{ executionDetails.id }}</span>
            </div>
            <div class="info-item">
              <strong>项目ID:</strong>
              <span>{{ executionDetails.project_id }}</span>
            </div>
            <div class="info-item">
              <strong>状态:</strong>
              <span>{{ executionDetails.status }}</span>
            </div>
            <div class="info-item">
              <strong>平台:</strong>
              <span>{{ executionDetails.platform }}</span>
            </div>
            <div class="info-item">
              <strong>开始时间:</strong>
              <span>{{ executionDetails.start_time }}</span>
            </div>
            <div class="info-item">
              <strong>结束时间:</strong>
              <span>{{ executionDetails.end_time }}</span>
            </div>
            <div class="info-item">
              <strong>持续时间:</strong>
              <span>{{ executionDetails.duration }}秒</span>
            </div>
          </div>

          <!-- 步骤执行详情 -->
          <div class="steps-section">
            <h3>步骤执行详情</h3>
            <div v-if="executionDetails.logs && executionDetails.logs.length > 0">
              <div class="stage-container" v-for="(stageLogs, stage) in groupedLogs" :key="stage">
                <div class="stage-header">
                  <h4>{{ stage }} 阶段</h4>
                </div>
                <div class="stage-content">
                  <div class="step-container" v-for="(stepLogs, step) in groupLogsByStep(stageLogs)" :key="step">
                    <div class="step-header">
                      <h5>{{ step }}</h5>
                    </div>
                    <div class="step-content">
                      <div class="log-entry" v-for="log in stepLogs" :key="log.id">
                        <div class="log-time">{{ log.timestamp }}</div>
                        <div class="log-level" :class="log.level">{{ log.level }}</div>
                        <div class="log-message">{{ log.message }}</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else>
              <p>暂无执行日志</p>
            </div>
          </div>

          <!-- 执行指标 -->
          <div class="metrics-section" v-if="executionDetails.metrics">
            <h3>执行指标</h3>
            <div class="metrics-grid">
              <div class="metric-card">
                <h4>总执行时间</h4>
                <div class="value">{{ executionDetails.metrics.total_duration }}s</div>
              </div>
              <div class="metric-card">
                <h4>CPU使用率</h4>
                <div class="value">{{ executionDetails.metrics.cpu_usage || 0 }}%</div>
              </div>
              <div class="metric-card">
                <h4>内存使用率</h4>
                <div class="value">{{ executionDetails.metrics.memory_usage || 0 }}%</div>
              </div>
              <div class="metric-card">
                <h4>测试覆盖率</h4>
                <div class="value">{{ executionDetails.metrics.test_coverage || 0 }}%</div>
              </div>
            </div>
          </div>
        </div>
        <div v-else>
          <p>加载中...</p>
        </div>
      </section>

      <!-- 模板管理 -->
      <section v-else-if="activeTab === 'templates'">
        <h2>模板管理</h2>
        
        <!-- 创建模板表单 -->
        <div class="form-section">
          <h3>{{ editingTemplate ? '编辑模板' : '创建新模板' }}</h3>
          <div class="form-group">
            <label for="templatePlatform">平台</label>
            <select id="templatePlatform" v-model="newTemplate.platform">
              <option value="mock">Mock (模拟平台)</option>
              <option value="github_actions">GitHub Actions</option>
            </select>
          </div>
          <div class="form-group">
            <label for="templateLanguage">语言</label>
            <input type="text" id="templateLanguage" v-model="newTemplate.language" placeholder="输入语言，例如：Go">
          </div>
          <div class="form-group">
            <label for="templateFramework">框架</label>
            <input type="text" id="templateFramework" v-model="newTemplate.framework" placeholder="输入框架，例如：React">
          </div>
          <div class="form-group">
            <label for="templateFilename">文件名</label>
            <input type="text" id="templateFilename" v-model="newTemplate.filename" placeholder="输入文件名，例如：.github/workflows/ci.yml">
          </div>
          <div class="form-group">
            <label for="templateConfigType">配置类型</label>
            <select id="templateConfigType" v-model="newTemplate.configType">
              <option value="yaml">YAML</option>
              <option value="json">JSON</option>
            </select>
          </div>
          <div class="form-group">
            <label for="templateContent">模板内容</label>
            <textarea id="templateContent" v-model="newTemplate.content" placeholder="输入模板内容" rows="10"></textarea>
          </div>
          <div class="form-actions">
            <button class="secondary" @click="cancelEditTemplate">取消</button>
            <button class="primary" @click="saveTemplate">{{ editingTemplate ? '更新模板' : '创建模板' }}</button>
          </div>
        </div>

        <!-- 模板列表 -->
        <div class="template-list">
          <h3>模板列表</h3>
          <div v-if="loadingTemplates">加载中...</div>
          <div v-else-if="templates.length === 0">暂无模板</div>
          <div v-else>
            <div v-for="template in templates" :key="template.id" class="template-item">
              <div class="template-header">
                <h4>{{ template.platform }} - {{ template.language || '通用' }}</h4>
                <span v-if="template.is_builtin" class="builtin-badge">内置</span>
              </div>
              <p><strong>文件名:</strong> {{ template.filename }}</p>
              <p><strong>配置类型:</strong> {{ template.config_type }}</p>
              <div class="template-actions">
                <button class="secondary" @click="editTemplate(template)">编辑</button>
                <button v-if="!template.is_builtin" class="danger" @click="deleteTemplate(template.id)">删除</button>
                <button v-if="template.is_builtin" class="secondary" @click="resetTemplate(template.id)">重置</button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 使用指南 -->
      <section v-else-if="activeTab === 'guide'">
        <h2>使用指南</h2>
        
        <!-- 如何使用前端页面 -->
        <div class="guide-section">
          <h3>如何使用前端页面</h3>
          <div class="guide-content">
            <h4>1. 项目管理</h4>
            <p><strong>创建项目:</strong> 在"项目管理"页面，填写项目名称和路径，点击"创建项目"按钮。</p>
            <p><strong>分析技术栈:</strong> 在项目列表中，找到要分析的项目，点击"分析技术栈"按钮。</p>
            <p><strong>生成管道配置:</strong> 在项目列表中，找到要生成配置的项目，点击"生成管道"按钮，选择平台和模板，点击"确定"按钮。</p>
            <p><strong>执行管道:</strong> 在项目列表中，找到要执行的项目，点击"执行管道"按钮，选择平台，点击"确定"按钮。</p>
            <p><strong>查看执行历史:</strong> 在项目列表中，找到要查看的项目，点击"查看执行"按钮。</p>
            <p><strong>分析优化建议:</strong> 在项目列表中，找到要分析的项目，点击"分析优化"按钮。</p>
          </div>
          
          <div class="guide-content">
            <h4>2. 模板管理</h4>
            <p><strong>创建模板:</strong> 在"模板管理"页面，填写模板信息，点击"创建模板"按钮。</p>
            <p><strong>编辑模板:</strong> 在模板列表中，找到要编辑的模板，点击"编辑"按钮，修改模板信息，点击"更新模板"按钮。</p>
            <p><strong>删除模板:</strong> 在模板列表中，找到要删除的模板，点击"删除"按钮，确认删除操作。</p>
            <p><strong>重置内置模板:</strong> 在模板列表中，找到要重置的内置模板，点击"重置"按钮。</p>
          </div>
          
          <div class="guide-content">
            <h4>3. 导航操作</h4>
            <p><strong>回到上一页:</strong> 点击导航栏左侧的"← 回到上一页"按钮。</p>
            <p><strong>切换选项卡:</strong> 点击导航栏中的选项卡名称，切换到对应页面。</p>
          </div>
        </div>
        
        <!-- 设计文档 -->
        <div class="guide-section">
          <h3>设计文档</h3>
          <div class="guide-content">
            <h4>1. 系统架构</h4>
            <p><strong>后端架构:</strong> 使用Go语言开发，SQLite数据库，RESTful API接口。</p>
            <p><strong>前端架构:</strong> 使用Vue 3框架，Vite构建工具，Axios HTTP客户端。</p>
            <p><strong>核心模块:</strong> 项目管理、技术栈分析、CI/CD配置生成、模板管理、执行管理、指标分析、优化建议。</p>
          </div>
          
          <div class="guide-content">
            <h4>2. 核心功能</h4>
            <p><strong>技术栈分析:</strong> 自动识别项目的技术栈，包括语言、框架、构建工具等。</p>
            <p><strong>CI/CD配置生成:</strong> 根据技术栈生成适合的CI/CD配置文件，支持GitHub Actions和Mock平台。</p>
            <p><strong>模板管理:</strong> 管理CI/CD配置模板，支持内置模板和自定义模板。</p>
            <p><strong>执行管理:</strong> 执行CI/CD管道，查看执行历史和指标。</p>
            <p><strong>优化建议:</strong> 分析执行数据，生成优化建议。</p>
          </div>
          
          <div class="guide-content">
            <h4>3. 技术实现</h4>
            <p><strong>模板系统:</strong> 支持内置模板和自定义模板，根据技术栈和平台选择最适合的模板。</p>
            <p><strong>执行系统:</strong> 支持Mock平台的完整执行功能，用于测试和开发；GitHub平台仅支持CI配置生成，不支持完整执行功能。</p>
            <p><strong>Mock平台意义:</strong> Mock平台的主要作用是测试整个CI/CD流程功能，包括配置生成、执行流程、指标收集和优化建议等，无需实际连接GitHub或其他CI/CD平台，便于开发和测试阶段使用。</p>
            <p><strong>指标系统:</strong> 收集执行时长、成功率、CPU/内存使用率等指标。</p>
            <p><strong>优化系统:</strong> 分析历史执行数据，识别瓶颈，生成优化建议。</p>
          </div>
          
          <div class="guide-content">
            <h4>4. API接口</h4>
            
            <!-- 项目管理 API -->
            <div class="api-category">
              <h5>项目管理</h5>
              <ul class="api-list">
                <li><code>GET /api/v1/projects</code> - 获取项目列表</li>
                <li><code>POST /api/v1/projects</code> - 创建新项目</li>
                <li><code>GET /api/v1/projects/{id}</code> - 获取项目详情</li>
                <li><code>PUT /api/v1/projects/{id}</code> - 更新项目</li>
                <li><code>DELETE /api/v1/projects/{id}</code> - 删除项目</li>
              </ul>
            </div>
            
            <!-- 技术栈分析 API -->
            <div class="api-category">
              <h5>技术栈分析</h5>
              <ul class="api-list">
                <li><code>POST /api/v1/projects/{id}/analyze</code> - 分析项目技术栈</li>
                <li><code>GET /api/v1/projects/{id}/tech-stack</code> - 获取技术栈分析结果</li>
              </ul>
            </div>
            
            <!-- CI/CD 配置生成 API -->
            <div class="api-category">
              <h5>CI/CD 配置生成</h5>
              <ul class="api-list">
                <li><code>POST /api/v1/projects/{id}/generate-pipeline</code> - 生成管道配置</li>
              </ul>
            </div>
            
            <!-- 模板管理 API -->
            <div class="api-category">
              <h5>模板管理</h5>
              <ul class="api-list">
                <li><code>GET /api/v1/templates</code> - 获取模板列表</li>
                <li><code>POST /api/v1/templates</code> - 创建新模板</li>
                <li><code>GET /api/v1/templates/{id}</code> - 获取模板详情</li>
                <li><code>PUT /api/v1/templates/{id}</code> - 更新模板</li>
                <li><code>DELETE /api/v1/templates/{id}</code> - 删除模板</li>
                <li><code>POST /api/v1/templates/{id}/reset</code> - 重置内置模板</li>
              </ul>
            </div>
            
            <!-- 执行管理 API -->
            <div class="api-category">
              <h5>执行管理</h5>
              <ul class="api-list">
                <li><code>POST /api/v1/projects/{id}/execute</code> - 执行管道</li>
                <li><code>GET /api/v1/projects/{id}/executions</code> - 获取执行历史</li>
                <li><code>GET /api/v1/executions/{id}</code> - 获取执行详情</li>
                <li><code>POST /api/v1/executions/{id}/stop</code> - 停止执行</li>
                <li><code>GET /api/v1/executions/{id}/metrics</code> - 获取执行指标</li>
                <li><code>GET /api/v1/executions/{id}/logs</code> - 获取执行日志</li>
              </ul>
            </div>
            
            <!-- 指标分析 API -->
            <div class="api-category">
              <h5>指标分析</h5>
              <ul class="api-list">
                <li><code>GET /api/v1/projects/{id}/metrics</code> - 获取项目指标</li>
              </ul>
            </div>
            
            <!-- 优化建议 API -->
            <div class="api-category">
              <h5>优化建议</h5>
              <ul class="api-list">
                <li><code>POST /api/v1/projects/{id}/analyze-optimization</code> - 分析优化建议</li>
                <li><code>GET /api/v1/projects/{id}/optimization-suggestions</code> - 获取优化建议</li>
              </ul>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- 消息提示 -->
    <div v-if="message" :class="['message', message.type]">
      {{ message.text }}
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',
  data() {
    return {
      activeTab: 'projects',
      projects: [],
      loading: false,
      newProject: {
        name: '',
        path: ''
      },
      techStackResult: null,
      pipelineResult: null,
      executions: [],
      projectMetrics: null,
      optimizationResult: null,
      executionDetails: null,
      message: null,
      // 平台选择相关
      showPlatformDialog: false,
      selectedPlatform: 'mock',
      selectedTemplateId: 0,
      platformActionType: '', // 'generate' 或 'execute'
      currentProjectId: null,
      // 轮询相关
      pollingInterval: null,
      pollingProjectId: null,
      // 模板管理相关
      templates: [],
      loadingTemplates: false,
      newTemplate: {
        platform: 'mock',
        language: '',
        framework: '',
        filename: '',
        configType: 'yaml',
        content: ''
      },
      editingTemplate: null,
      currentTemplateId: null
    };
  },
  mounted() {
    this.loadProjects();
  },
  watch: {
    // 监听标签页变化，切换到模板管理页面时加载模板列表
    activeTab(newTab) {
      if (newTab === 'templates') {
        this.loadTemplates();
      }
    }
  },
  computed: {
    // 按阶段分组日志
    groupedLogs() {
      if (!this.executionDetails || !this.executionDetails.logs) {
        return {};
      }
      
      const grouped = {};
      this.executionDetails.logs.forEach(log => {
        if (!grouped[log.stage]) {
          grouped[log.stage] = [];
        }
        grouped[log.stage].push(log);
      });
      return grouped;
    }
  },
  methods: {
    // 按步骤分组日志
    groupLogsByStep(logs) {
      if (!logs || logs.length === 0) {
        return {};
      }
      
      const grouped = {};
      logs.forEach(log => {
        const step = log.step || 'default';
        if (!grouped[step]) {
          grouped[step] = [];
        }
        grouped[step].push(log);
      });
      return grouped;
    },

    // 加载项目列表
    async loadProjects() {
      this.loading = true;
      try {
        const response = await axios.get('/api/v1/projects');
        this.projects = response.data.data || [];
      } catch (error) {
        this.showMessage('加载项目失败', 'error');
        console.error('Error loading projects:', error);
      } finally {
        this.loading = false;
      }
    },

    // 创建项目
    async createProject() {
      if (!this.newProject.name || !this.newProject.path) {
        this.showMessage('请填写项目名称和路径', 'error');
        return;
      }

      try {
        const response = await axios.post('/api/v1/projects', this.newProject);
        if (response.data.status === 'success' && response.data.data) {
          // 保存路径信息到项目对象中
          const projectWithPath = {
            ...response.data.data,
            path: this.newProject.path
          };
          this.projects.push(projectWithPath);
          this.newProject = { name: '', path: '' };
          this.showMessage('项目创建成功', 'success');
        } else {
          this.showMessage('创建项目失败: ' + (response.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage('创建项目失败', 'error');
        console.error('Error creating project:', error);
      }
    },

    // 删除项目
    async deleteProject(projectId) {
      if (confirm('确定要删除这个项目吗？')) {
        try {
          const response = await axios.delete(`/api/v1/projects/${projectId}`);
          if (response.data.status === 'success') {
            this.projects = this.projects.filter(p => p.id !== projectId);
            this.showMessage('项目删除成功', 'success');
          } else {
            this.showMessage('删除项目失败: ' + (response.data.message || '未知错误'), 'error');
          }
        } catch (error) {
          this.showMessage('删除项目失败', 'error');
          console.error('Error deleting project:', error);
        }
      }
    },

    // 分析技术栈
    async analyzeTechStack(projectId) {
      try {
        const analyzeResponse = await axios.post(`/api/v1/projects/${projectId}/analyze`);
        if (analyzeResponse.data.status === 'success') {
          const response = await axios.get(`/api/v1/projects/${projectId}/tech-stack`);
          this.techStackResult = response.data.data || {};
          this.activeTab = 'tech-stack';
          this.showMessage('技术栈分析成功', 'success');
        } else {
          this.showMessage('技术栈分析失败: ' + (analyzeResponse.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage('技术栈分析失败', 'error');
        console.error('Error analyzing tech stack:', error);
      }
    },

    // 生成管道配置
    async generatePipeline(projectId) {
      try {
        // 先获取项目信息，以获取项目路径
        const projectResponse = await axios.get(`/api/v1/projects/${projectId}`);
        if (projectResponse.data.status === 'success' && projectResponse.data.data) {
          const projectPath = projectResponse.data.data.path;
          const response = await axios.post(`/api/v1/projects/${projectId}/generate-pipeline?platform=mock&path=${encodeURIComponent(projectPath)}`);
          if (response.data.status === 'success') {
            this.pipelineResult = response.data.data || '';
            this.activeTab = 'pipeline';
            this.showMessage('管道配置生成成功', 'success');
          } else {
            this.showMessage('管道配置生成失败: ' + (response.data.message || '未知错误'), 'error');
          }
        } else {
          this.showMessage('获取项目信息失败', 'error');
        }
      } catch (error) {
        this.showMessage('管道配置生成失败', 'error');
        console.error('Error generating pipeline:', error);
      }
    },

    // 执行管道
    async executePipeline(projectId) {
      try {
        const response = await axios.post(`/api/v1/projects/${projectId}/execute?platform=mock`);
        if (response.data.status === 'success') {
          this.showMessage('管道执行已启动', 'success');
          // 加载执行历史
          this.viewExecutions(projectId);
        } else {
          this.showMessage('管道执行失败: ' + (response.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage('管道执行失败', 'error');
        console.error('Error executing pipeline:', error);
      }
    },

    // 查看执行历史
    async viewExecutions(projectId) {
      try {
        const response = await axios.get(`/api/v1/projects/${projectId}/executions`);
        this.executions = response.data.data || [];
        this.activeTab = 'execution';
        this.pollingProjectId = projectId;
        this.startPolling();
      } catch (error) {
        this.showMessage('加载执行历史失败', 'error');
        console.error('Error loading executions:', error);
      }
    },

    // 启动轮询
    startPolling() {
      // 先停止之前的轮询
      this.stopPolling();
      
      // 每3秒轮询一次
      this.pollingInterval = setInterval(() => {
        this.pollExecutions();
      }, 3000);
    },

    // 停止轮询
    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null;
      }
    },

    // 轮询执行历史
    async pollExecutions() {
      if (!this.pollingProjectId || this.activeTab !== 'execution') {
        this.stopPolling();
        return;
      }
      
      try {
        const response = await axios.get(`/api/v1/projects/${this.pollingProjectId}/executions`);
        this.executions = response.data.data || [];
      } catch (error) {
        console.error('Error polling executions:', error);
      }
    },

    // 查看执行指标
    async getExecutionMetrics(executionId) {
      try {
        const response = await axios.get(`/api/v1/executions/${executionId}/metrics`);
        alert(JSON.stringify(response.data.data || {}, null, 2));
      } catch (error) {
        this.showMessage('获取执行指标失败', 'error');
        console.error('Error getting execution metrics:', error);
      }
    },

    // 查看执行日志
    async getExecutionLogs(executionId) {
      try {
        const response = await axios.get(`/api/v1/executions/${executionId}/logs`);
        alert(JSON.stringify(response.data.data || {}, null, 2));
      } catch (error) {
        this.showMessage('获取执行日志失败', 'error');
        console.error('Error getting execution logs:', error);
      }
    },

    // 停止执行
    async stopExecution(executionId) {
      try {
        const response = await axios.post(`/api/v1/executions/${executionId}/stop`);
        if (response.data.status === 'success') {
          this.showMessage('执行已停止', 'success');
          // 刷新执行列表
          this.loadProjects();
        } else {
          this.showMessage('停止执行失败: ' + (response.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage('停止执行失败', 'error');
        console.error('Error stopping execution:', error);
      }
    },

    // 查看执行详情
    async viewExecutionDetails(executionId) {
      try {
        const response = await axios.get(`/api/v1/executions/${executionId}`);
        if (response.data.status === 'success') {
          this.executionDetails = response.data.data;
          this.activeTab = 'execution-details';
          this.showMessage('获取执行详情成功', 'success');
        } else {
          this.showMessage('获取执行详情失败: ' + (response.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage('获取执行详情失败', 'error');
        console.error('Error getting execution details:', error);
      }
    },

    // 分析优化建议
    async analyzeOptimization(projectId) {
      try {
        const response = await axios.post(`/api/v1/projects/${projectId}/analyze-optimization`);
        if (response.data.status === 'success') {
          this.optimizationResult = response.data.data || {};
          this.activeTab = 'optimization';
          this.showMessage('优化分析成功', 'success');
        } else {
          this.showMessage('优化分析失败: ' + (response.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage('优化分析失败', 'error');
        console.error('Error analyzing optimization:', error);
      }
    },

    // 显示消息
    showMessage(text, type = 'info') {
      this.message = { text, type };
      setTimeout(() => {
        this.message = null;
      }, 3000);
    },

    // 显示平台选择对话框
    showPlatformSelector(actionType, projectId) {
      this.platformActionType = actionType;
      this.currentProjectId = projectId;
      this.selectedPlatform = 'mock'; // 默认选择 mock 平台
      this.selectedTemplateId = 0; // 默认选择默认模板
      this.loadTemplates(); // 加载模板列表
      this.showPlatformDialog = true;
    },

    // 确认平台选择
    confirmPlatformSelection() {
      if (!this.currentProjectId) {
        this.showMessage('项目ID不能为空', 'error');
        return;
      }

      if (this.platformActionType === 'generate') {
        this.generatePipelineWithPlatform(this.currentProjectId, this.selectedPlatform, this.selectedTemplateId);
      } else if (this.platformActionType === 'execute') {
        this.executePipelineWithPlatform(this.currentProjectId, this.selectedPlatform);
      }

      this.showPlatformDialog = false;
    },

    // 生成管道配置（带平台参数）
    async generatePipelineWithPlatform(projectId, platform, templateId = 0) {
      try {
        // 从前端保存的项目列表中获取项目路径
        const project = this.projects.find(p => p.id === projectId);
        if (!project) {
          this.showMessage('项目不存在', 'error');
          return;
        }
        
        const projectPath = project.path;
        if (!projectPath) {
          this.showMessage('项目路径为空', 'error');
          return;
        }
        
        let url = `/api/v1/projects/${projectId}/generate-pipeline?platform=${platform}&path=${encodeURIComponent(projectPath)}`;
        if (templateId > 0) {
          url += `&template_id=${templateId}`;
        }
        const response = await axios.post(url);
        if (response.data.status === 'success') {
          this.pipelineResult = response.data.data || '';
          this.activeTab = 'pipeline';
          this.showMessage('管道配置生成成功', 'success');
        } else {
          this.showMessage('管道配置生成失败: ' + (response.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage('管道配置生成失败', 'error');
        console.error('Error generating pipeline:', error);
      }
    },

    // 执行管道（带平台参数）
    async executePipelineWithPlatform(projectId, platform) {
      try {
        const response = await axios.post(`/api/v1/projects/${projectId}/execute?platform=${platform}`);
        if (response.data.status === 'success') {
          this.showMessage('管道执行已启动', 'success');
          // 加载执行历史
          this.viewExecutions(projectId);
        } else {
          this.showMessage('管道执行失败: ' + (response.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage('管道执行失败', 'error');
        console.error('Error executing pipeline:', error);
      }
    },

    // 加载模板列表
    async loadTemplates() {
      this.loadingTemplates = true;
      try {
        const response = await axios.get('/api/v1/templates');
        if (response.data.status === 'success') {
          this.templates = response.data.data || [];
        } else {
          this.showMessage('加载模板列表失败: ' + (response.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage('加载模板列表失败', 'error');
        console.error('Error loading templates:', error);
      } finally {
        this.loadingTemplates = false;
      }
    },

    // 创建或更新模板
    async saveTemplate() {
      try {
        if (!this.newTemplate.platform || !this.newTemplate.filename || !this.newTemplate.configType || !this.newTemplate.content) {
          this.showMessage('平台、文件名、配置类型和内容为必填字段', 'error');
          return;
        }

        const templateData = {
          platform: this.newTemplate.platform,
          language: this.newTemplate.language,
          framework: this.newTemplate.framework,
          filename: this.newTemplate.filename,
          config_type: this.newTemplate.configType,
          content: this.newTemplate.content
        };

        let response;
        if (this.editingTemplate) {
          // 更新模板
          response = await axios.put(`/api/v1/templates/${this.currentTemplateId}`, templateData);
        } else {
          // 创建模板
          response = await axios.post('/api/v1/templates', templateData);
        }

        if (response.data.status === 'success') {
          this.showMessage(this.editingTemplate ? '模板更新成功' : '模板创建成功', 'success');
          this.cancelEditTemplate();
          this.loadTemplates();
        } else {
          this.showMessage((this.editingTemplate ? '模板更新失败' : '模板创建失败') + ': ' + (response.data.message || '未知错误'), 'error');
        }
      } catch (error) {
        this.showMessage(this.editingTemplate ? '模板更新失败' : '模板创建失败', 'error');
        console.error('Error saving template:', error);
      }
    },

    // 编辑模板
    editTemplate(template) {
      this.editingTemplate = true;
      this.currentTemplateId = template.id;
      this.newTemplate = {
        platform: template.platform,
        language: template.language,
        framework: template.framework,
        filename: template.filename,
        configType: template.config_type,
        content: template.content
      };
    },

    // 取消编辑模板
    cancelEditTemplate() {
      this.editingTemplate = false;
      this.currentTemplateId = null;
      this.newTemplate = {
        platform: 'mock',
        language: '',
        framework: '',
        filename: '',
        configType: 'yaml',
        content: ''
      };
    },

    // 删除模板
    async deleteTemplate(templateId) {
      if (confirm('确定要删除这个模板吗？')) {
        try {
          const response = await axios.delete(`/api/v1/templates/${templateId}`);
          if (response.data.status === 'success') {
            this.showMessage('模板删除成功', 'success');
            this.loadTemplates();
          } else {
            this.showMessage('删除模板失败: ' + (response.data.message || '未知错误'), 'error');
          }
        } catch (error) {
          this.showMessage('删除模板失败', 'error');
          console.error('Error deleting template:', error);
        }
      }
    },

    // 重置模板
    async resetTemplate(templateId) {
      if (confirm('确定要重置这个内置模板吗？')) {
        try {
          const response = await axios.post(`/api/v1/templates/${templateId}/reset`);
          if (response.data.status === 'success') {
            this.showMessage('模板重置成功', 'success');
            this.loadTemplates();
          } else {
            this.showMessage('重置模板失败: ' + (response.data.message || '未知错误'), 'error');
          }
        } catch (error) {
          this.showMessage('重置模板失败', 'error');
          console.error('Error resetting template:', error);
        }
      }
    },

    // 回到上一页
    goBack() {
      window.history.back();
    }
  }
};
</script>

<style scoped>
.app {
  min-height: 100vh;
}

.form-section {
  margin-bottom: 30px;
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #ddd;
}

.form-section h3 {
  margin-bottom: 15px;
  color: #4a90e2;
}

.message {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 20px;
  border-radius: 4px;
  color: white;
  font-weight: 500;
  z-index: 1000;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.message.success {
  background-color: #28a745;
}

.message.error {
  background-color: #dc3545;
}

.message.info {
  background-color: #17a2b8;
}

pre {
  background-color: #f5f5f5;
  padding: 15px;
  border-radius: 4px;
  overflow-x: auto;
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
  line-height: 1.4;
  margin-top: 10px;
}

/* 平台选择对话框样式 */
.platform-dialog {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.dialog-content {
  background-color: white;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  max-width: 400px;
  width: 90%;
}

.dialog-content h3 {
  margin-bottom: 20px;
  color: #4a90e2;
}

.dialog-content p {
  margin-bottom: 20px;
  color: #666;
}

.platform-options {
  margin-bottom: 25px;
}

.platform-options label {
  display: block;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.platform-options label:hover {
  background-color: #f5f5f5;
  border-color: #4a90e2;
}

.platform-options input[type="radio"] {
  margin-right: 10px;
}

.dialog-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

/* 模板管理相关样式 */
.template-list {
  margin-top: 30px;
}

.template-item {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #ddd;
}

.template-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.template-header h4 {
  margin: 0;
  color: #4a90e2;
}

.builtin-badge {
  background-color: #28a745;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.template-actions {
  margin-top: 15px;
  display: flex;
  gap: 10px;
}

.form-actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
}

textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
  resize: vertical;
}

select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
  font-size: 14px;
}

/* 回到上一页按钮样式 */
.back-button {
  background-color: #6c757d;
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.back-button:hover {
  background-color: #5a6268;
}

.back-button:active {
  background-color: #495057;
}

/* 使用指南相关样式 */
.guide-section {
  margin-bottom: 40px;
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #ddd;
}

.guide-section h3 {
  margin-bottom: 20px;
  color: #4a90e2;
  border-bottom: 1px solid #ddd;
  padding-bottom: 10px;
}

.guide-content {
  margin-bottom: 25px;
  padding-left: 20px;
}

.guide-content h4 {
  margin-bottom: 15px;
  color: #333;
  font-size: 16px;
}

.guide-content p {
  margin-bottom: 10px;
  line-height: 1.5;
  color: #666;
}

.guide-content strong {
  color: #333;
  font-weight: 500;
}

/* API接口相关样式 */
.api-category {
  margin-bottom: 25px;
  padding-left: 20px;
}

.api-category h5 {
  margin-bottom: 12px;
  color: #4a90e2;
  font-size: 14px;
  font-weight: 600;
}

.api-list {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.api-list li {
  margin-bottom: 8px;
  padding-left: 20px;
  position: relative;
}

.api-list li:before {
  content: "•";
  color: #4a90e2;
  font-weight: bold;
  position: absolute;
  left: 0;
}

.api-list code {
  background-color: #f5f5f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  color: #d73a49;
  margin-right: 8px;
}

/* 技术栈分析相关样式 */
.tech-stack-details {
  background-color: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #ddd;
  padding: 20px;
}

.tech-stack-item {
  margin-bottom: 20px;
}

.tech-stack-item strong {
  display: block;
  margin-bottom: 8px;
  color: #4a90e2;
  font-size: 16px;
  font-weight: 600;
}

.tech-stack-item span {
  font-size: 14px;
  color: #333;
}

/* 依赖包列表样式 */
.dependencies-list {
  margin-top: 10px;
}

.dependencies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 10px;
}

.dependency-item {
  background-color: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  border-left: 4px solid #4a90e2;
}

.dependency-name {
  display: block;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  color: #333;
  margin-bottom: 4px;
}

.dependency-version {
  font-size: 12px;
  color: #666;
}

/* 文件列表样式 */
.files-list {
  margin-top: 10px;
}

.files-list ul {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.files-list li {
  margin-bottom: 8px;
  padding-left: 20px;
  position: relative;
  font-size: 14px;
  color: #333;
}

.files-list li:before {
  content: "📄";
  position: absolute;
  left: 0;
}

/* 执行详情页面样式 */
.execution-info {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
  margin-bottom: 30px;
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #ddd;
}

.info-item {
  display: flex;
  flex-direction: column;
}

.info-item strong {
  margin-bottom: 5px;
  color: #4a90e2;
  font-size: 14px;
}

.info-item span {
  font-size: 14px;
  color: #333;
}

.steps-section {
  margin-bottom: 30px;
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #ddd;
}

.steps-section h3 {
  margin-bottom: 20px;
  color: #4a90e2;
  border-bottom: 1px solid #ddd;
  padding-bottom: 10px;
}

.stage-container {
  margin-bottom: 25px;
}

.stage-header {
  background-color: #4a90e2;
  color: white;
  padding: 10px 15px;
  border-radius: 4px 4px 0 0;
}

.stage-header h4 {
  margin: 0;
  font-size: 16px;
}

.stage-content {
  background-color: white;
  border: 1px solid #ddd;
  border-top: none;
  border-radius: 0 0 4px 4px;
  padding: 15px;
}

.step-container {
  margin-bottom: 20px;
}

.step-header {
  background-color: #f0f8ff;
  border-left: 4px solid #4a90e2;
  padding: 8px 12px;
  margin-bottom: 10px;
}

.step-header h5 {
  margin: 0;
  font-size: 14px;
  color: #333;
}

.step-content {
  padding-left: 10px;
}

.log-entry {
  display: flex;
  align-items: flex-start;
  margin-bottom: 10px;
  padding: 8px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.log-time {
  width: 150px;
  font-size: 12px;
  color: #666;
  margin-right: 15px;
}

.log-level {
  width: 80px;
  font-size: 12px;
  font-weight: bold;
  margin-right: 15px;
  text-transform: uppercase;
}

.log-level.info {
  color: #17a2b8;
}

.log-level.error {
  color: #dc3545;
}

.log-level.warn {
  color: #ffc107;
}

.log-level.debug {
  color: #6c757d;
}

.log-message {
  flex: 1;
  font-size: 14px;
  color: #333;
  line-height: 1.4;
}

.metrics-section {
  margin-bottom: 30px;
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #ddd;
}

.metrics-section h3 {
  margin-bottom: 20px;
  color: #4a90e2;
  border-bottom: 1px solid #ddd;
  padding-bottom: 10px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.metric-card {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #ddd;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.metric-card h4 {
  margin-bottom: 10px;
  color: #666;
  font-size: 14px;
}

.metric-card .value {
  font-size: 24px;
  font-weight: bold;
  color: #4a90e2;
}

@media (max-width: 768px) {
  .dialog-content {
    padding: 20px;
  }
  
  .dialog-actions {
    flex-direction: column;
  }
  
  .dialog-actions button {
    width: 100%;
  }
  
  .template-actions {
    flex-direction: column;
  }
  
  .template-actions button {
    width: 100%;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .form-actions button {
    width: 100%;
  }
  
  .execution-info {
    grid-template-columns: 1fr;
  }
  
  .metrics-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .log-entry {
    flex-direction: column;
  }
  
  .log-time,
  .log-level {
    width: auto;
    margin-bottom: 5px;
  }
}
</style>
