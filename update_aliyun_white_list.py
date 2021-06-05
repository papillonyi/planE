import argh
import requests
import ipaddress

QUERY_IP_API = "http://members.3322.org/dyndns/getip"


def run():
    local_ip = get_local_ip()
    pass


def get_local_ip():
    response = requests.get(QUERY_IP_API)
    assert response.status_code == 200
    print(response.content)
    ip = ipaddress.ip_address(response.content)
    print(ip)
    return ip



if __name__ == "__main__":
    argh.dispatch_command(run)
