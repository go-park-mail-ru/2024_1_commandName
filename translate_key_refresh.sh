#bash
CONFIG_FILE="config.yml"
update_trKey(){
    trKey=$(yc iam create-token)
    bearer_token="Bearer $trKey"
    sed -i "s/trKey:.*/trKey: \"$bearer_token\"/" $CONFIG_FILE
    echo "trKey обновлен: $bearer_token"
}

# Бесконечный цикл для обновления каждые 10 часов
while true; do
    update_trKey
    sleep 36000 # 10 часов в секундах
done