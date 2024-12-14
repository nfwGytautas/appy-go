"""
Entry file for the query engine
"""

import glob
import subprocess
import argparse
import re
import pathlib

import drivers

from query import Query
from migration import Migration
from macro import Macro

from driver import DatabaseDriver


class QueryEngine:
    """
    The query function generator for
    """

    QUERY_MARKER = "@query"
    QUERY_END_MARKER = "@endquery"
    MACRO_MARKER = "@macro"
    MACRO_END_MARKER = "@endmacro"

    def __init__(self) -> None:
        arguments = argparse.ArgumentParser(description="Query Engine")

        arguments.add_argument(
            "query_dir", type=str, help="The directory of the queries")
        arguments.add_argument(
            "migration_dir", type=str, help="The directory of the migrations")
        arguments.add_argument("out_dir", type=str,
                               help="The output directory")
        arguments.add_argument(
            "package", type=str, help="The name of the package", default="package"
        )

        values = arguments.parse_args()

        self.__query_dir = values.query_dir
        self.__migration_dir = values.migration_dir
        self.__out_dir = values.out_dir
        self.__package = values.package

    def __get_queries(self) -> [str]:
        """
        Get all the query files in the query directory
        """
        return glob.glob(f"{self.__query_dir}/**/*.sql", recursive=True)

    def __get_migrations(self) -> [str]:
        """
        Get all the migration files in the migration directory
        """
        return glob.glob(f"{self.__migration_dir}/*.sql")

    def __parse_query_file(self, file: str) -> ([Query], [Macro]):
        macro_pattern = r"@macro\s(\S+)\n(.*?)@endmacro"
        query_pattern = r"@query\s(\S+)\n(.*?)@endquery"

        macro_query = re.compile(macro_pattern, re.DOTALL)
        query_query = re.compile(query_pattern, re.DOTALL)

        with open(file, 'r', encoding="utf-8") as f:
            query_file = f.read()

            macros = [Macro(x[0], x[1])
                      for x in macro_query.findall(query_file)]

            queries = [Query(str(pathlib.Path(file).relative_to(self.__query_dir)), x[0], x[1])
                       for x in query_query.findall(query_file)]

        return (queries, macros)

    def __parse_migration(self, origin: str) -> Migration:
        with open(origin, "r", encoding="utf-8") as f:
            return Migration(origin, f.read())

    def __get_driver(self) -> DatabaseDriver:
        # return drivers.MySqlDriver()
        return drivers.PostgresDriver()

    def __apply_macros(self, queries: [Query], macros: [Macro]) -> None:
        """
        Apply macros to the queries
        """
        for query in queries:
            query.set_content(
                self.__apply_macros_to_query(query.content, macros))

    def __apply_macros_to_query(self, content: str, macros: [Macro]) -> str:
        has_replaced = True

        while has_replaced:
            has_replaced = False
            for macro in macros:
                macro_key = f"%%{macro.name}%%"

                if macro_key not in content:
                    continue

                content = content.replace(
                    macro_key, macro.body)
                has_replaced = True

        return content

    def generate(self):
        """
        Generate the golang queries
        """

        # Load
        queries = []
        macros = []

        contents = [self.__parse_query_file(
            file) for file in self.__get_queries()]

        for content in contents:
            queries += content[0]
            macros += content[1]

        migrations = [self.__parse_migration(
            file) for file in self.__get_migrations()]

        print(
            f"Found: '{len(queries)}' queries from '{len(self.__get_queries())}' files")

        print(
            f"Found: '{len(macros)}' macros from '{len(self.__get_queries())}' files")

        print(
            f"Found: '{len(migrations)}' migrations from '{len(self.__get_migrations())}' files")

        # Generate the 'database.go'
        driver = self.__get_driver()

        print(f"Driver: '{driver.get_name()}'")

        # Apply macros
        self.__apply_macros(queries, macros)

        # Write base files
        driver.write_base(self.__package, self.__out_dir)

        # Generate queries
        driver.write_queries(self.__out_dir, queries)

        # Generate migrations
        driver.write_migrations(self.__out_dir, migrations)

        # Write constants
        driver.write_constants(self.__out_dir)

        # Tidy up all files
        print("Code generated, running 'gofmt' and 'goimports'")
        for file in glob.glob(f"{self.__out_dir}/*.go"):
            print(f"  - {file}")
            subprocess.run(["goimports", "-w", file], check=True)
            subprocess.run(["gofmt", "-w", file], check=True)

        print(f"Generating done, files stored in '{self.__out_dir}'")


if __name__ == "__main__":
    engine = QueryEngine()
    engine.generate()
