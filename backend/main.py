"""The app module, containing the app initialization."""

import logging
import os
from argparse import ArgumentParser

from rest.server import Server


def main():
    """An application factory."""

    def setup_logger(dbg):
        log_level = logging.INFO
        log_format = "%(asctime)s [%(levelname)s] %(message)s"

        if dbg:
            log_level = logging.DEBUG
            log_format = "%(asctime)s [%(levelname)s] %(filename)s: %(message)s"

        logging.basicConfig(level=log_level, format=log_format)

    # parsing cli flags
    parser = ArgumentParser(description="Image processing web-service")
    parser.add_argument(
        "--serviceurl",
        help="url of this web-service in format \"http://<addr>:<port>/\"",
        default=os.environ.get("SERVICEURL", None)
    )
    parser.add_argument(
        "--dbg",
        help="enable debug mode",
        default=os.environ.get("DEBUG", False)
    )

    args = parser.parse_args()
    if not args.serviceurl or len(args.serviceurl.split(':')) < 2:
        exit(parser.print_usage())

    host: str = args.serviceurl.split(':')[-2][2:]
    port: int = int(args.serviceurl.split(':')[-1][:-1])

    setup_logger(args.dbg)

    srv = Server(host=host, port=port, dbg=args.dbg)
    srv.run()


if __name__ == '__main__':
    main()
