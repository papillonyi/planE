import argh
import requests
import ipaddress
import logging
import json
from aliyunsdkcore.client import AcsClient
from aliyunsdkecs.request.v20140526.DescribeSecurityGroupAttributeRequest import (
    DescribeSecurityGroupAttributeRequest,
)
from aliyunsdkecs.request.v20140526.AuthorizeSecurityGroupRequest import (
    AuthorizeSecurityGroupRequest,
)

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(filename)s[line:%(lineno)d] %(levelname)s %(message)s",
    datefmt="%a, %d %b %Y %H:%M:%S",
)

QUERY_IP_API = "http://members.3322.org/dyndns/getip"


def get_aliyun_client(access_key_id, access_key_secret, region_id):
    client = AcsClient(access_key_id, access_key_secret, region_id)
    return client


def get_security_group(security_group_id, send_request_function):
    request = DescribeSecurityGroupAttributeRequest()
    request.set_SecurityGroupId(security_group_id)
    logging.info(f"describe security group {security_group_id} attribute")
    return send_request_function(request)


def revoke_security_group_rule():
    pass


def add_security_group_rule(
    home_permission, local_ip, security_group_id, send_request_function
):
    request = AuthorizeSecurityGroupRequest()
    request.set_Policy(home_permission["Policy"])
    request.set_Description(home_permission["Description"])
    request.set_Priority(home_permission["Priority"])
    request.set_NicType(home_permission["NicType"])
    request.set_PortRange(home_permission["PortRange"])
    request.set_SourceCidrIp(local_ip)
    request.set_IpProtocol(home_permission["IpProtocol"])
    request.set_SecurityGroupId(security_group_id)
    logging.info(
        f"'add local ip {local_ip} permission to security group {security_group_id}"
    )
    return send_request_function(request)


def get_send_request_function(client: AcsClient):
    def _send_request(request):
        """
        send open api request
        :param request:
        :return:
        """
        request.set_accept_format("json")
        try:
            response_str = client.do_action_with_exception(request)
            logging.debug(response_str)
            response_detail = json.loads(response_str)
            return response_detail
        except Exception as e:
            logging.error(e)

    return _send_request


def get_local_ip():
    response = requests.get(QUERY_IP_API)
    assert response.status_code == 200
    ip = ipaddress.ip_address(response.content.decode().replace("\n", ""))
    logging.info(f"Get local ip {ip}")
    return ip


def run(access_key_id, access_key_secret, security_group_id, region_id="cn-shanghai"):
    aliyun_client = get_aliyun_client(access_key_id, access_key_secret, region_id)
    send_request_function = get_send_request_function(aliyun_client)
    local_ip = str(get_local_ip())
    local_ip = "1.1.1.1"
    permissions = get_security_group(security_group_id, send_request_function)[
        "Permissions"
    ]["Permission"]
    home_permission = next(filter(lambda x: x["Description"] == "home", permissions))

    if local_ip == home_permission["SourceCidrIp"]:
        logging.info(
            f"local ip {local_ip} equal to source cider ip {home_permission['SourceCidrIp']}"
        )
    else:
        logging.info(
            f"local ip {local_ip} not equal to source cider ip {home_permission['SourceCidrIp']}"
        )
        add_security_group_rule(
            home_permission, local_ip, security_group_id, send_request_function
        )


if __name__ == "__main__":
    argh.dispatch_command(run)
