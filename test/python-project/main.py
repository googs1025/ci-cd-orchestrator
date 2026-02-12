#!/usr/bin/env python3

print("Hello, Python Project!")

# 简单的 Django 项目检查
print("\nChecking Django installation...")
try:
    import django
    print(f"Django version: {django.__version__}")
except ImportError:
    print("Django not installed")

# 简单的 pytest 检查
print("\nChecking pytest installation...")
try:
    import pytest
    print(f"pytest version: {pytest.__version__}")
except ImportError:
    print("pytest not installed")
