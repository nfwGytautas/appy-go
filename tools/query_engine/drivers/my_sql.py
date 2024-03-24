from query import Query
from driver import DatabaseDriver
import sys
sys.path.append("..")


class MySqlDriver(DatabaseDriver):
    """
    MySql implementation
    """

    def get_name(self) -> str:
        return "MySql"

    def _write_constants(self, file: any) -> None:
        file.write('const p_CONTEXT_DRIVER = "mysql"\n')

    def _format_query(self, query: str) -> str:
        return f'"{Query.ARGUMENT_PATTERN.sub("?", query)}"'

    def _format_migration(self, query: str) -> str:
        return f'"{query}"'
