import winreg

# Функция проверки наличия sub_key System (по умолчанию его обычно нет).
def system():
    key_to_read = r'Software\Microsoft\Windows\CurrentVersion\Policies\System'
    try:
        reg = ConnectRegistry(None, HKEY_CURRENT_USER)
        k = OpenKey(reg, key_to_read)
        print("Есть")     # Действие если есть.
    except:
        print("Добавить") # Действие если нет (например добавить).

# Функция проверки наличия ключа блокировки (невозможности открытия) диспетчера задач.
def proverka_registry():
    value = [] # Массив, куда будут записываться все ключи sub_key System.
    j = 0
    startup = reg.OpenKey(reg.HKEY_CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\System")
    while 1:
        try:
            value.append(reg.EnumValue(startup, j)) # Запись всех ключей sub_key System в массив.
            j += 1
        except:
            break
    print(value)
    if (('DisableTaskMgr', 1, 4)) in value: # Проверка массива на ключ.
        print("OK")              # Действие, если ключ есть.
    else:
        print("Что-нибудь")      # Действие, если ключа нет (например добавление оного).

# Функция проверки наличия ключа блокировки (невозможности открытия) встроенного редактора реестра.       
def proverka_dispetcher():
    value = [] # Массив, куда будут записываться все ключи sub_key System.
    j = 0
    startup = reg.OpenKey(reg.HKEY_CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\System")
    while 1:
        try:
            value.append(reg.EnumValue(startup, j)) # Запись всех ключей sub_key System в массив.
            j += 1
        except:
            break
    print(value)
    if (('DisableRegistryTools', 1, 4)) in value: # Проверка массива на ключ.
        print("OK")              # Действие, если ключ есть.
    else:
        print("Что-нибудь")      # Действие, если ключа нет (например добавление оного) .       
