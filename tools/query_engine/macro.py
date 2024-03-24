
class Macro:
    """
    A @macro statement
    """

    def __init__(self, name, body):
        self.__name = name
        self.__body = body

    @property
    def name(self):
        return self.__name

    @property
    def body(self):
        return ' '.join([x.rstrip() for x in self.__body.replace('\t', ' ').split()])

    def __repr__(self) -> str:
        return self.__str__()

    def __str__(self) -> str:
        return f"Macro: [Name: '{self.__name}']"
