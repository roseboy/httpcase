#setup.py

from setuptools import setup

setup(
    name = "httpcase",
    version = "1.0.17",
    author = "Mr.K",
    author_email = "roseboy@live.com",
    description = ("HttpCase - api auto test tool."),
    url = "https://github.com/roseboy/httpcase",
    install_requires = [
        'requests>=2.19.1',
        'tqdm>=4.61.1'
    ],
    packages=['httpcase'],
    entry_points={
        'console_scripts': ['hc=httpcase.httpcase:main'],
    }
)