FROM python:3.8
RUN mkdir ~/pip
ADD pip.conf ~/pip/
ADD requirements.txt ./
RUN pip install -r  requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple
ADD update_aliyun_white_list.py ./

CMD python ./update_aliyun_white_list.py $access_key_id $access_key_secret $security_group_id