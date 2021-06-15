#!/usr/bin/python
# -*- coding: UTF-8 -*-
import os
import stat
import sys
import platform
import tempfile
import tarfile
import zipfile
import re
import requests

home_page = "https://github.com/roseboy/httpcase"

def main():
    
    urls= get_latest_release_url()
    print(urls)
    if len(urls) == 0:
        print("Download HttpCase error")
        exit(1)
    download_url = get_my_os_url(urls)
    print(download_url)
    if download_url == "":
        print("can't get download url for your os")
        exit(1)

    print("Download:" + download_url)
    r = requests.get(download_url)
    file_name = tempfile.gettempdir() + os.sep + download_url[download_url.rindex("/") + 1:]
    with open(file_name, "wb") as code:
        code.write(r.content)
    code.close()

    print("DownloadComplete:" + file_name)

    installer = Installer(file_name)
    installer.install("hc")

    os.remove(file_name)
    return

def get_latest_release_url():
    git_version=""
    release_urls = []
    response  = requests.get(home_page + "/releases")

    tags = re.findall( r'<span class="css-truncate-target"(.*?)</span>', response.text)
    if len(tags) > 0:
        git_version = tags[0]
        git_version = git_version[git_version.find(">")+1:]
    else:
        return release_urls

    urls = re.findall( r'<a href="(.*?)"', response.text)
    for url in urls:
        if url.find(git_version)>-1 and url.find("download")>-1 and not url.endswith(".txt"):
            release_urls.append(url)

    return release_urls

def get_my_os_url(urls):
    h = Hardware()
    os_tag = ""
    if h.is_mac_os():
        os_tag = "darwin_x86_64"
    elif h.is_linux() and h.is_intel_cpu() and h.is_64bit_cpu():
        os_tag = "linux_x86_64"
    elif h.is_linux() and h.is_arm_cpu() and h.is_64bit_cpu():
        os_tag = "linux_armv64"
    elif h.is_linux() and h.is_arm_cpu() and not h.is_64bit_cpu():
        os_tag = "linux_armv6"
    elif h.is_windows() and h.is_intel_cpu and h.is_64bit_cpu:
        os_tag = "windows_x86_64"
    elif h.is_windows() and h.is_intel_cpu() and not h.is_64bit_cpu:
        os_tag = "windows_i386"
    elif h.is_windows() and h.is_arm_cpu():
        os_tag = "windows_armv6"
    if os_tag=="":
        return ""
    for url in urls:
        if url.find(os_tag) >-1:
            return "https://github.com"+url
    return ""
        

class Installer:
    hardware = ""
    tar = ""
    file = ""
    fileDir = ""
    binDir = ""

    def __init__(self, file):
        self.hardware = Hardware()
        self.tar = Tar(file)
        self.file = file
        self.fileDir = "httpcase"
        self.binDir = sys.path[0]

    def install(self, bin_name):
        if self.hardware.is_mac_os():
            self.mac_install(bin_name)
        elif self.hardware.is_linux():
            self.linux_install(bin_name)
        elif self.hardware.is_windows():
            self.windows_install(bin_name)

    def windows_install(self, bin_name):
        if self.binDir.endswith(".exe"):
            self.binDir=self.binDir[:self.binDir.rindex(os.sep)]

        target_dir = os.environ['LOCALAPPDATA'] + os.sep + "Programs" + os.sep +self.fileDir
        target_bin = self.binDir + os.sep + bin_name + ".bat"
        origin_bin = self.binDir + os.sep + bin_name + ".exe"

        self.tar.extract(target_dir)
        with open(target_bin, 'w') as wf:
            wf.write("@"+target_dir+os.sep+bin_name+".exe %1 %2 %3 %4 %5 %6 %7 %8 %9")
        wf.close()

        with open(origin_bin+".bat", 'w') as wf:
            wf.write("@choice /t 1 /d y /n >nul\r\n")
            wf.write("@del "+origin_bin+"\r\n")
            wf.write("@del "+origin_bin+".bat"+"\r\n")
        wf.close()
        os.system("start /min \"cmd /c "+origin_bin+".bat\"")


    def linux_install(self, bin_name):
        target_dir = "/usr/local/" + self.fileDir
        target_bin = "/usr/local/bin/" + bin_name
        origin_bin = self.binDir + os.sep + bin_name
        if not os.access(target_dir, os.W_OK):
            target_dir = self.binDir + os.sep + ".." + os.sep + self.fileDir
            target_bin = origin_bin
        
        self.tar.extract(target_dir)

        if os.path.exists(origin_bin):
            os.remove(origin_bin)
        os.symlink(target_dir+os.sep+bin_name, target_bin)
        
        # os.chmod(target_bin, stat.S_IRWXU | stat.S_IRWXG | stat.S_IRWXO)
        # os.chmod(origin_bin, stat.S_IRWXU | stat.S_IRWXG | stat.S_IRWXO)

    def mac_install(self, bin_name):
        self.linux_install(bin_name)


class Hardware:
    system = ""
    machine = ""
    architecture = ""

    def __init__(self):
        uname = platform.uname()
        if type(uname) == tuple:
            self.system = uname[0]
            self.machine = uname[4]
            self.architecture = platform.architecture()[0]
        else:
            self.system = uname.system
            self.machine = uname.machine
            self.architecture = platform.architecture()[0]

    def is_mac_os(self):
        return self.system.lower() == "darwin"

    def is_linux(self):
        return self.system.lower() == "linux"

    def is_windows(self):
        return self.system.lower() == "windows"

    def is_intel_cpu(self):
        mc = self.machine.lower()
        return mc == "amd64" or mc == "x86_64" or mc == "i386"

    def is_arm_cpu(self):
        return self.machine.lower().find("arm") > -1

    def is_64bit_cpu(self):
        return self.architecture.find("64") > -1


class Tar:
    file = ""

    def __init__(self, file):
        self.file = file

    def extract(self, target):
        if self.file.endswith(".zip"):
            self.un_zip(target)
        elif self.file.endswith(".tar.gz") or self.file.endswith(".tgz"):
            self.un_tgz(target)
        else:
            print("Error:not support")

    def un_tgz(self, target):
        tar = tarfile.open(self.file)
        if os.path.isdir(target):
            pass
        else:
            os.makedirs(target)
        for name in tar.getnames():
            tar.extract(name, target)
        tar.close()

    def un_zip(self, target):
        zip_file = zipfile.ZipFile(self.file)
        if os.path.isdir(target):
            pass
        else:
            os.makedirs(target)
        for names in zip_file.namelist():
            zip_file.extract(names, target)
        zip_file.close()


if __name__ == '__main__':
    main()

# uname_result(system='Linux', node='VM-0-9-ubuntu', release='5.4.0-72-generic', version='#80-Ubuntu SMP Mon Apr 12 17:35:00 UTC 2021', machine='x86_64', processor='x86_64')
# uname_result(system='Windows', node='DESKTOP-6BV4V7I', release='10', version='10.0.17763', machine='AMD64', processor='Intel64 Family 6 Model 60 Stepping 3, GenuineIntel')
# uname_result(system='Darwin', node='MrKdeMacBook-Pro.local', release='20.5.0', version='Darwin Kernel Version 20.5.0: Sat May  8 05:10:33 PDT 2021; root:xnu-7195.121.3~9/RELEASE_X86_64', machine='x86_64', processor='i386')
# uname_result(system='Linux', node='raspberrypi', release='5.10.11-v7+', version='#1399 SMP Thu Jan 28 12:06:05 GMT 2021', machine='armv7l', processor='')
# uname_result(system='Linux', node='raspberrypi', release='5.10.17-v7l+', version='#1403 SMP Mon Feb 22 11:33:35 GMT 2021', machine='armv7l', processor='')
