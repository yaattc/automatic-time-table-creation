"""Rest package defines everything related to REST API."""
from flask import Flask


class Server:
    """Server specifies the router, registers blueprints and initializes the REST API overall."""
    def __init__(self, host: str, port: int, dbg: bool):
        """Initialize the server parameters."""
        self.__flask = Flask(__name__)
        self.__host = host
        self.__port = port
        self.__dbg = dbg

    def run(self):
        """Run the server."""
        self.__flask.run(host=self.__host, port=self.__port, debug=self.__dbg, use_reloader=False)
