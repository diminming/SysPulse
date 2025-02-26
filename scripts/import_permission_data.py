import json

def read_json_from_file():
    with open('/workspace/SysPulse/scripts/permission_data.json') as f:
        return json.load(f)

#[
#   {
#     "id": 1,
#     "identity": "create_biz",
#     "method": "POST",
#     "url": "/api/biz",
#     "createTimestamp": 1740364669000,
#     "updateTimestamp": 1740364669000
#   },
#   {
#     "id": 2,
#     "identity": "update_biz",
#     "method": "PUT",
#     "url": "/api/biz",
#     "createTimestamp": 1740364669000,
#     "updateTimestamp": 1740364669000
#   },
#]
def import_permission_data():
    data = read_json_from_file()
    for item in data:
        sql = f"INSERT INTO permission (`id`, `identity`, `name`, `method`, `url`, `createTimestamp`, `updateTimestamp`) VALUES (%d, '%s', '%s', '%s', '%s', %d, %d);" % (item['id'], item['identity'], item['name'], item['method'], item['url'], item['createTimestamp'], item['updateTimestamp'])
        print(sql)

if __name__ == "__main__":
    import_permission_data()