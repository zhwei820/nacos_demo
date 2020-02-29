import time
import nacos
from main import config
from anotherfile import print_conf as print_conf_in_anotherfile
from imp import reload


SERVER_ADDRESSES = "localhost:8848"

client = nacos.NacosClient(SERVER_ADDRESSES, )

# get config
d = "natsmicro_dev"
g = "natsmicro"
print(client.get_config(d, g))


class Share:
    content = None
    count = 0


def test_cb(args):
    Share.count += 1
    print(client.get_config(d, g))
    print(config.conf)
    config.reload_conf()
    reload(config)
    print(config.conf)
    print_conf_in_anotherfile()

    print(Share.count)


# client.add_config_watcher(d, g, test_cb)
# client.add_config_watcher(d, g, test_cb)
client.add_config_watcher(d, g, test_cb)
time.sleep(1000)
