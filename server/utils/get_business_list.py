# -*- coding: utf-8 -*-
import os
import sys
import json
import subprocess
from zipfile import ZipFile


# 执行代码块返回代码执行中环境变量
def exec_script(code, filename='<script>', optimize=-1, global_vars=None):
    """
    :param code: 创建代码对象的源。 这可以是字符串，字节字符串或AST模块对象
    :param filename: 代码文件名称，如果不是从文件读取代码则传递一些可辨认的值
    :param global_vars: 执行代码的全局变量
    :param optimize: 编译器的优化级别
    :return: 成功返回(True, global_vars)  失败返回(False, str(e))
    """
    if global_vars is None:
        global_vars = {}
    try:
        code_obj = compile(source=code, filename=filename, mode='exec', optimize=optimize)
        exec(code_obj, global_vars)
        result = True, global_vars
    except Exception as e:
        print(e)
        result = False, str(e)
    return result


# 从python代码中获取特定的key-value
def obtain_parm(code, keys: list = None):
    """
    :param code: python代码字符串，字节字符串或AST模块对象
    :param keys: 需要提取的key列表
    :return: dict
    """
    parm_dict = {}
    executed, result = exec_script(code)

    if executed:
        for key in keys:
            if key in result:
                parm_dict[key] = result.get(key)
    return parm_dict


# .ipynb转py
def ipynb2py(ipynb_file, py_file):
    """
    :param ipynb_file: .ipynb文件路径
    :param py_file: .py文件路径
    :return:  True 成功  False 失败
    """
    CMD = "/cmbp/Python370/bin/ipynb-py-convert {} {}".format(ipynb_file, py_file)
    p = subprocess.Popen(CMD, stderr=subprocess.DEVNULL, stdout=subprocess.DEVNULL, shell=True,
                         executable='bash')
    p.wait()
    if p.returncode != 0:
        return False
    return True


# 生成config.py
def make_config_py(dir_path):
    """
    :param dir_path: 目录路径
    :return: True 成功  False 失败
    """
    ipynb_file = os.path.join(dir_path, "config.ipynb")
    py_file = os.path.join(dir_path, "config.py")

    if os.path.exists(py_file) and os.path.exists(ipynb_file):
        if os.path.getmtime(ipynb_file) > os.path.getmtime(py_file):
            return ipynb2py(ipynb_file, py_file)
    elif os.path.exists(ipynb_file):
        return ipynb2py(ipynb_file, py_file)
    elif os.path.exists(py_file):
        return True
    else:
        return False
    return True


# 变量类型转换
def typeof(variate):
    var_type = "string"
    if isinstance(variate, int):
        var_type = "int"
    elif isinstance(variate, str):
        var_type = "string"
    elif isinstance(variate, float):
        var_type = "float"
    elif isinstance(variate, list):
        var_type = "array"
    elif isinstance(variate, dict):
        var_type = "object"
    return var_type


# 构造business_list
def make_business_list(parm_dict):
    """
    :param parm_dict:
    {"config": {"config_id": "test",
                "location": "test",
                "json_url": "http://192.168.188.106:1112",
                "business_params": {
                    "一般危险区域": {
                        "area": [[[0, 0], [1280, 0], [1280, 720], [0, 720]]],
                        "threshold": 0.5
                    }
                }},
     "business_type": {
         "一般危险区域": ["person", "helmet"]
     }
     }
     }
    :return: list
    [{
    "business_name": "",
    "business_params": [
        {"default": "0", "nullable": False, "param_name": "target", "param_type": "int", "paras_desc": "target",
         "paras_text": "target"}],
    "business_type": []}]
    """
    business_list = []
    if parm_dict:
        business_type = parm_dict.get("business_type")
        config = parm_dict.get("config")
        if config:
            business_params = config.get("business_params")
            if business_params:
                for k, v in business_params.items():
                    business = {}
                    business["business_name"] = k
                    business["business_params"] = []
                    for i_key, i_value in v.items():
                        var_type = typeof(i_value)
                        i_dict = {
                            "default": json.dumps(i_value, ensure_ascii=False) if var_type != 'string' else i_value,
                            "nullable": False,
                            "param_name": i_key,
                            "param_type": var_type,
                            "paras_desc": i_key,
                            "paras_text": i_key,
                            'param_shape': 0
                        }
                        business["business_params"].append(i_dict)
                    if business_type and k in business_type:
                        business["business_type"] = business_type[k]
                    else:
                        business["business_type"] = []
                    business_list.append(business)
        elif business_type:
            for k, v in business_type.items():
                business = {}
                business["business_name"] = k
                business["business_type"] = v
                business["business_params"] = []
                business_list.append(business)
    return business_list


# 从python文件中获取特定的key-value
def obtain_parm_by_path(file_path, keys: list = None):
    """
    :param file_path: 要读取的python文件路径
    :param keys: 需要提取的key列表
    :return: list
    """
    try:
        with open(file_path, 'r') as f:
            code = f.read()
    except Exception as e:
        print(e)
        return dict()
    parm_dict = obtain_parm(code, keys)

    return make_business_list(parm_dict)


# 从zip文件中获取业务config的key-value
def obtain_parm_in_zip(zip_file_path, keys: list = None, ignore_dir_name: list=None):
    """
    :param ignore_dir_name: list plugins下要忽略的文件夹 default为[".idea", "__pycache__"]
    :param zip_file_path: 要读取的zip文件路径
    :param keys: 需要提取的key列表
    :return: list
    """
    # 从zip压缩包中获取特定业务config.py的代码
    if ignore_dir_name is None:
        ignore_dir_name = [".idea", "__pycache__"]
    config_codes = []
    with ZipFile(zip_file_path, 'r') as f_zip:
        config_list = []
        for f_name in f_zip.namelist():
            file_name_list = f_name.split('/')
            if len(file_name_list) == 3 and (
                    file_name_list[0] == "plugins" and
                    file_name_list[1] not in ignore_dir_name and
                    file_name_list[2] == "config.py"):
                config_list.append(f_name)
        for zip_config_path in config_list:
            config_codes.append(f_zip.read(zip_config_path))

    # 从config文件中提取特定字典
    business_list = []
    if config_codes:
        for config_code in config_codes:
            parm_dict = obtain_parm(config_code, keys)
            business_list += make_business_list(parm_dict)
    return business_list


if __name__ == '__main__':
    zip_obj = sys.argv[1]
    if zip_obj:
        res = obtain_parm_in_zip(zip_obj, keys=['config', 'business_type'])
#         print(res)
#         print(type(res))
#         for i in res:
#             print(type(i))
        print(json.dumps(res))
#         print(type(json.dumps(res)))