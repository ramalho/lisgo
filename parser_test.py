from lis import Expression, parse, lispstr

from pytest import mark

@mark.parametrize(
    'source, expected',
    [
        ('7', 7),
        ('x', 'x'),
        ('(sum 1 2 3)', ['sum', 1, 2, 3]),
        ('(+ (* 2 100) (* 1 10))', ['+', ['*', 2, 100], ['*', 1, 10]]),
        ('99 100', 99),  # parse stops at the first complete expression
        ('(a)(b)', ['a']),
    ],
)
def test_parse(source: str, expected: Expression) -> None:
    got = parse(source)
    assert got == expected



@mark.parametrize(
    'obj, expected',
    [
        (0, '0'),
        (1, '1'),
        (False, '#f'),
        (True, '#t'),
        (1.5, '1.5'),
        ('sin', 'sin'),
        (['+', 1, 2], '(+ 1 2)'),
        (['if', ['<', 'a', 'b'], True, False], '(if (< a b) #t #f)'),
        ([], '()'),
        (None, 'None'),
        (..., 'Ellipsis'),
    ],
)
def test_lispstr(obj: object, expected: str) -> None:
    got = lispstr(obj)
    assert got == expected