import pathlib
import re

class Migration:
    """
    Migration class
    """

    REQUIRES_PATTERN = re.compile(r'@Requires:\s*([^ \t\n\r\f\v]+)')

    def __init__(self, origin: str, content :str) -> None:
        self.__origin = origin
        self.__content = content

    @property
    def name(self) -> str:
        return pathlib.Path(self.__origin).stem

    def get_content(self) -> str:
        """
        Get pre-processed content of the migration
        """

        # Remove comments
        lines = self.__content.split("\n")
        lines = [x.lstrip().rstrip() for x in lines if not x.lstrip().startswith("--")]

        return "\\n".join(lines)

    def get_dependancies(self) -> str:
        """
        Find dependancy of the migration
        """
        matches = self.REQUIRES_PATTERN.findall(self.__content)

        if len(matches) > 1:
            raise ValueError("More than 1 @Requires is not allowed")

        if len(matches) == 1:
            return matches[0]

        return ""
