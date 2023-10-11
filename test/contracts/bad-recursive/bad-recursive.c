typedef unsigned char byte;
typedef unsigned int i32;
typedef unsigned long long i64;
typedef unsigned int bigInt;

bigInt bigIntNew(long long value);
void bigIntFinishUnsigned(bigInt reference);

void init()
{
}

i64 doStackoverflow(i64 a)
{
    if (a % 2 == 0)
    {
        return 42;
    }

    i64 x = doStackoverflow(a * 8 + 1);
    i64 y = doStackoverflow(a * 2 + 1);
    return x + y + a;
}

void badRecursive()
{
    i64 result = doStackoverflow(1);
    bigInt resultBig = bigIntNew(result);
    bigIntFinishUnsigned(resultBig);
}
