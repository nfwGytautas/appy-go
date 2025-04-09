"""
Database driver interface
"""

import pathlib
import shutil
from typing import List
from functools import reduce

from query import Query, QueryResult
from macro import Macro
from migration import Migration


class DatabaseDriver:
    """
    Interface for datbase driver
    """

    GO_VERSION = "1.21.3"

    INITIALIZE_CODE = """
        import (
            appy_driver "github.com/nfwGytautas/appy-go/driver"
            appy_logger "github.com/nfwGytautas/appy-go/logger"
        )

        func Initialize(connectionString string) error {
            return appy_driver.Initialize(appy_driver.InitializeArgs{
                ConnectionString: connectionString,
                Version:          DATAMODEL_VERSION,
                Migration:        g_ROOT_MIGRATION,
            })
        }
    """

    CONSTANTS_CODE = """
        import (
        \tappy_driver "github.com/nfwGytautas/appy-go/driver"
        \tappy "github.com/nfwGytautas/appy"
        )
    """

    def __init__(self) -> None:
        self.__depend_list = []
        self.__migration_list = []
        self.__package_name = ""

    def get_name(self) -> str:
        """
        Get the name of the driver
        """
        raise NotImplementedError("Not implemented")

    def write_base(self, package_name: str, out_dir: str) -> None:
        """
        Write base files to the directory
        """
        self.__package_name = package_name

        self.__create_root(out_dir)


        with open(out_dir + "constants.go", "w+", encoding="utf-8") as f:
            f.write(f"package {package_name}_driver\n")
            f.write(DatabaseDriver.CONSTANTS_CODE)

        with open(out_dir + "driver.go", "w+", encoding="utf-8") as f:
            f.write(f"package {package_name}_driver\n")
            f.write(DatabaseDriver.INITIALIZE_CODE)

    def write_constants(self, out_dir: str) -> None:
        """
        Writes constants.go
        """
        migration_root = list(
            set(self.__migration_list) - set(self.__depend_list))

        if len(migration_root) > 1:
            raise ValueError(f"More than 1 final migration '{migration_root}'")

        if len(migration_root) == 0:
            raise ValueError("No final migration")

        with open(out_dir + "/constants.go", "a", encoding="utf-8") as f:
            f.write("\n")
            self._write_constants(f)
            f.write(f'const DATAMODEL_VERSION = "{migration_root[0]}"\n')
            f.write(
                f'var g_ROOT_MIGRATION appy_driver.MigrationFn = m{migration_root[0]}\n')

    def write_queries(self, out_dir: str, queries: List[Query]) -> None:
        """
        Write queries to the out directory
        """

        print("Generating queries:")

        for i, q in enumerate(queries):
            print(f"  - [{i+1:3}/{len(queries):<3}] {q.name}")
            file_path = f"{out_dir}/{q.origin_escaped}.go"

            if not pathlib.Path(file_path).exists():
                self.__write_implementation_template(file_path)

            self.__write_query(file_path, q)

    def write_migrations(self, out_dir: str, migrations: List[Migration]) -> None:
        """
        Generates migrations
        """

        print("Generating migrations")

        for i, m in enumerate(migrations):
            print(f"  - [{i+1:3}/{len(migrations):<3}] {m.name}")

            self.__write_migration(f"{out_dir}/{m.name}.go", m)

    def _write_constants(self, file: any) -> None:
        raise NotImplementedError("Not impemented")

    def _format_query(self, query: str) -> str:
        raise NotImplementedError("Not implemented")

    def _format_migration(self, query: str) -> str:
        raise NotImplementedError("Not implemented")

    def __create_root(self, out_dir: str) -> None:
        p = pathlib.Path(out_dir)

        if p.exists():
            print("Removing old directory")
            if not p.is_dir():
                raise RuntimeError("Out exists but is not a directory")

            shutil.rmtree(out_dir)

        p.mkdir(parents=True, exist_ok=False)

    def __write_implementation_template(self, file: str) -> None:
        with open(file, "w+", encoding="utf-8") as f:
            f.write(f"package {self.__package_name}_driver\n")
            f.write("\n")
            f.write("import (\n")
            f.write('\t"database/sql"\n')
            f.write('\t"github.com/jackc/pgx/v5/pgtype"\n')
            f.write('\tappy_driver "github.com/nfwGytautas/appy-go/driver"\n')
            f.write('\tappy_logger "github.com/nfwGytautas/appy-go/logger"\n')
            f.write(")\n\n")

    def __write_query(self, file: str, query: Query) -> None:
        args = query.get_args()
        result = query.get_result()
        seen_args = []
        args_string = ""

        with open(file, "a+", encoding="utf-8") as f:
            f.write(f"// Origin: '{query.origin}'\n")
            f.write(f"// Arguments: '{len(args)}'\n")
            f.write(f"// Return: '{result.type}'\n")

            # Signature
            f.write(f"func Q{query.name}(tx *appy_driver.Tx")

            for arg in args:
                if arg.name not in seen_args:
                    f.write(f", {arg.name} {arg.type}")
                    seen_args.append(arg.name)

                if not arg.is_array:
                    args_string = args_string + arg.name + ", "
                else:
                    args_string = args_string + f"pq.Array({arg.name})" + ", "

            args_string = args_string[:-2]

            f.write(") ")

            # Return type
            f.write(
                f"({self.__get_result_equivalent(result.type)}, error)")

            # Function
            f.write("{\n")
            f.write(f'\tconst query = {self._format_query(query.content)}\n')

            # Logging
            f.write(f'\tappy_logger.Logger().Info("Executing \'{query.name}\'")\n')

            self.__write_implementation(f, result, args_string)

            f.write("}\n\n")

    def __get_result_equivalent(self, atype: str) -> str:
        if atype == "Row":
            return "appy_driver.RowResult"
        if atype == "Rows":
            return "appy_driver.RowsResult"
        if atype == "Null":
            return "appy_driver.ExecResult"

        return "ERROR"

    def __write_implementation(self, f: any, result: QueryResult, args: str) -> None:
        if result.type == "Row":
            return self.__row_return(f, args)
        if result.type == "Rows":
            return self.__rows_return(f, args)

        self.__nil_return(f, args)

    def __row_return(self, f: any, args: str) -> None:
        f.write("\trow := tx.QueryRow(query")
        if args != "":
            f.write(f", {args}")

        f.write(")\n")
        f.write("\treturn row, row.Err()\n")

    def __rows_return(self, f: any, args: str) -> None:
        f.write("\treturn tx.Query(query")

        if args != "":
            f.write(f", {args}")

        f.write(")\n")

    def __nil_return(self, f: any, args: str) -> None:
        f.write(f"\treturn tx.Exec(query, {args})\n")

    def __write_migration(self, f: str, migration: Migration) -> None:
        deps = migration.get_dependancies()

        self.__migration_list.append(migration.name)

        with open(f, "a+", encoding="utf-8") as f:
            f.write(f"package {self.__package_name}_driver\n")
            f.write("\n")
            f.write("import (\n")
            f.write('\t"database/sql"\n')
            f.write('\t"github.com/jackc/pgx/v5/pgtype"\n')
            f.write('\tappy_driver "github.com/nfwGytautas/appy-go/driver"\n')
            f.write('\tappy_logger "github.com/nfwGytautas/appy-go/logger"\n')
            f.write(")\n\n")

            # Function
            f.write(
                f"func m{migration.name}(tx *appy_driver.Tx, currentVersion string) error ")
            f.write("{\n")

            # Query constant
            f.write(
                f'\tconst query = {self._format_migration(migration.get_content())}\n')
            f.write('\tvar err error\n')

            # Version check
            f.write(f'\n\tif currentVersion == "{migration.name}" ')
            f.write("{\n")
            f.write('\t\tappy_logger.Logger().Info("    - Already at this version")\n')
            f.write("\t\treturn nil\n")
            f.write("\t}\n")

            # Dependancy
            if deps != "":
                # Logging
                f.write(f'\n\tappy_logger.Logger().Info("  - Dependency to \'{deps}\'")\n')

                f.write(f"\n\terr = m{deps}(tx, currentVersion)\n")
                f.write("\tif err != nil {\n")
                f.write("\t\treturn err\n")
                f.write("\t}\n")

                self.__depend_list.append(deps)

            # Logging
            f.write(
                f'\n\tappy_logger.Logger().Info("  - Migrating to \'{migration.name}\'")\n\n')

            # Execute migration
            f.write(
                f'\terr = appy_driver.MigrateToVersion(tx, "{migration.name}", query)\n')
            f.write('\tif err != nil {')
            f.write('\t\treturn err\n')
            f.write('\t}\n\n')
            f.write(f'\tcurrentVersion = "{migration.name}"\n')
            f.write('\t\nreturn nil\n')
            f.write("}\n\n")
