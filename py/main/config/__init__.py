
import yaml
import os


def get_yaml_data(file_path):
    try:
        with open(file_path, 'r') as f:
            file_data = f.read()
            data = yaml.load(file_data, Loader=yaml.FullLoader)
            return data
    except Exception as e:
        return {}

current_path = os.path.abspath(".")
yaml_path = os.path.join(
    current_path, "nacos-data/snapshot/natsmicro_dev+natsmicro+")

conf = get_yaml_data(yaml_path)


def reload_conf():
    global conf
    conf = get_yaml_data(yaml_path)