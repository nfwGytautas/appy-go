import pathlib
import re
from typing import List


class QueryArg:
    """
    Query argument class
    """

    def __init__(self, atype: str, name: str) -> None:
        self.__type = atype
        self.__name = name

    @property
    def name(self) -> str:
        return self.__name

    @property
    def type(self) -> str:
        return self.__type

    @property
    def is_array(self) -> bool:
        return self.__type.startswith("[]")


class QueryResult:
    """
    Query result class
    """

    def __init__(self, atype: str) -> None:
        self.__type = atype

    @property
    def type(self) -> str:
        return self.__type

    @property
    def is_null(self) -> bool:
        return self.__type == "Null"


class Query:
    """
    Single SQL query
    """

    ARGUMENT_PATTERN = re.compile(r'\$(.*?)\:(\w+)')

    def __init__(self, origin: str, name: str, content: str) -> None:
        self.__name = name
        self.__content = content
        self.__origin = origin

    def __repr__(self) -> str:
        return self.__str__()

    def __str__(self) -> str:
        return f"Query: [Origin: '{self.__origin}', Name: '{self.__name}']"

    @property
    def origin_escaped(self) -> str:
        return str(pathlib.Path('_'.join(self.__origin.split("/"))).stem)

    @property
    def origin(self) -> str:
        return self.__origin

    @property
    def name(self) -> str:
        return self.__name.split(":")[0]

    @property
    def content(self) -> str:
        return ' '.join([x.rstrip() for x in self.__content.replace('\t', ' ').split()])

    def set_content(self, content: str) -> None:
        self.__content = content

    def get_args(self) -> List[QueryArg]:
        """
        Get arguments of the query

        Argument format:
            $<type>:<name>
        Example:
            $uint64:index
        """
        matches = self.ARGUMENT_PATTERN.findall(self.__content)
        return [QueryArg(match[0], match[1]) for match in matches]

    def get_result(self) -> QueryResult:
        """
        Get query result

        Format:
            @query <name>:<result>
        Example:
            @query Query:Rows
            @query Query:Row
            @query Query:Null or @query Query
        """
        split = self.__name.split(":")

        if len(split) == 1:
            # No result 'Null' implied
            return QueryResult("Null")

        return QueryResult(split[1])
