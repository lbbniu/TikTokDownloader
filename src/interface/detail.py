from typing import Callable
from typing import TYPE_CHECKING
from typing import Union

from src.interface.template import API
from src.testers import Params

if TYPE_CHECKING:
    from src.config import Parameter


class Detail(API):
    def __init__(self,
                 params: Union["Parameter", Params],
                 cookie: str = None,
                 proxy: str = None,
                 detail_id: str = ...,
                 ):
        print("Detail __init__", "headers = ", params.headers, "params = ", self.params)
        super().__init__(params, cookie, proxy, )
        self.detail_id = detail_id
        self.api = f"{self.domain}aweme/v1/web/aweme/detail/"
        self.text = "作品数据"

    def generate_params(self, ) -> dict:
        return {'device_platform': 'webapp', 'aid': '6383', 'channel': 'channel_pc_web', 'update_version_code': '170400', 'pc_client_type': '1', 'version_code': '190500', 'version_name': '19.5.0', 'cookie_enabled': 'true', 'screen_width': '1536', 'screen_height': '864', 'browser_language': 'zh-SG', 'browser_platform': 'Win32', 'browser_name': 'Chrome', 'browser_version': '126.0.0.0', 'browser_online': 'true', 'engine_name': 'Blink', 'engine_version': '126.0.0.0', 'os_name': 'Windows', 'os_version': '10', 'cpu_core_num': '16', 'device_memory': '8', 'platform': 'PC', 'downlink': '10', 'effective_type': '4g', 'round_trip_time': '200', 'msToken': 'eHUQHQOZgTUdIyobTzkIBOxmCGDUmm6PTJzDi2PtXcP5XHCEKVrdcCNcfE8DhShYk_1P3llPBA6BYia8HNE7HcSMdpuV_XFOURF9gbEHnwolgwUzy9j12lL1UYekBA==', 'aweme_id': '6870423037087436046'}

    async def run(self,
                  referer: str = None,
                  single_page=True,
                  data_key: str = "aweme_detail",
                  error_text="",
                  cursor="cursor",
                  has_more="has_more",
                  params: Callable = lambda: {},
                  data: Callable = lambda: {},
                  method="GET",
                  headers: dict = None,
                  *args,
                  **kwargs,
                  ):
        print(f"作品 {self.detail_id} 获取数据", "params", params(), "data", data())
        return await super().run(
            referer,
            single_page,
            data_key,
            error_text or f"作品 {self.detail_id} 获取数据失败",
            cursor,
            has_more,
            params,
            data,
            method,
            headers,
            *args,
            **kwargs,
        )

    def check_response(self,
                       data_dict: dict,
                       data_key: str,
                       error_text="",
                       cursor="cursor",
                       has_more="has_more",
                       *args,
                       **kwargs,
                       ):
        try:
            if not (d := data_dict[data_key]):
                self.log.info(error_text)
            else:
                self.response = d
        except KeyError:
            self.log.error(f"数据解析失败，请告知作者处理: {data_dict}")
