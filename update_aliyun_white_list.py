import argh
import requests
import ipaddress
import logging
from aliyunsdkcore.client import AcsClient

QUERY_IP_API = "http://members.3322.org/dyndns/getip"
logger = logging.getLogger("update aliyun white list")
logger.setLevel(logging.DEBUG)
log_handler = logging.StreamHandler()
log_handler.setLevel(logging.INFO)
formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
log_handler.setFormatter(formatter)
logger.addHandler(log_handler)


def get_aliyun_client(access_key_id, access_key_secret, region_id):
    client = AcsClient(access_key_id, access_key_secret, region_id)
    return client


def run(access_key_id, access_key_secret, region_id="cn-shanghai"):
    aliyun_client = get_aliyun_client(access_key_id, access_key_secret, region_id)
    local_ip = get_local_ip()
    pass


def get_local_ip():
    response = requests.get(QUERY_IP_API)
    assert response.status_code == 200
    ip = ipaddress.ip_address(response.content.decode().replace("\n", ""))
    logger.info(f"Get local ip {ip}")
    return ip


if __name__ == "__main__":
    argh.dispatch_command(run)
