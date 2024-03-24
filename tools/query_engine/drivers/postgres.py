import re
from query import Query
from driver import DatabaseDriver
import sys
sys.path.append("..")


class PostgresDriver(DatabaseDriver):
    """
    Postgres implementation
    """

    def __init__(self) -> None:
        super().__init__()

        self.__index = -1

    def get_name(self) -> str:
        return "Postgres"

    def _write_constants(self, file: any) -> None:
        file.write('const p_CONTEXT_DRIVER = "postgres"\n')

    def _format_query(self, query: str) -> str:
        self.__index = 0
        def replace_with_increased_index(match):
            self.__index += 1
            return f"${self.__index}"

        query = self.__flatten(query)
        query = re.sub(Query.ARGUMENT_PATTERN, replace_with_increased_index, query)

        return f"`{query}`"

    def _format_migration(self, query: str) -> str:
        query = self.__flatten(query)
        return f'`{query}`'

    def __flatten(self, query: str) -> str:
        return " ".join(query.split("\\n"))
