import sys
import requests
import re

def get_product_info(product_url):
    # Извлекаем ID товара
    match = re.search(r'/catalog/(\d+)/detail', product_url)
    if not match:
        return "Ссылка некорректна. Убедись, что она ведёт на страницу товара."

    product_id = match.group(1)
    api_url = f'https://card.wb.ru/cards/v1/detail?appType=1&curr=rub&dest=-1257786&nm={product_id}'
    headers = {
        "User-Agent": "Mozilla/5.0"
    }

    response = requests.get(api_url, headers=headers)
    if response.status_code != 200:
        return f"Ошибка при запросе: {response.status_code}"

    try:
        product = response.json()['data']['products'][0]
        name = product['name']
        # price = product['priceU'] // 100
        sale_price = product['salePriceU'] // 100
        return name, sale_price
    except (KeyError, IndexError):
        return "Не удалось извлечь информацию. Возможно, товар не существует или временно недоступен."

# Пример запуска
if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Не передана ссылка.")
        exit(1)
    url = sys.argv[1]
    print(get_product_info(url))