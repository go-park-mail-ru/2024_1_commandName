#!/bin/bash

# Путь к файлу конфигурации
CONFIG_FILE="ProjectMessenger/cmd/messenger/config.yml"

# Функция для обновления trKey
update_trKey() {
    trKey=$(yc iam create-token)
    # Обновляем значение trKey в файле конфигурации
    sed -i "s/trKey:.*/trKey: \"$trKey\"/" $CONFIG_FILE
    echo "trKey обновлен: $trKey"
}

# Бесконечный цикл для обновления каждые 10 часов
while true; do
    update_trKey
    sleep 36000 # 10 часов в секундах
done