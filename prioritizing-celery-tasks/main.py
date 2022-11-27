from tasks import a, b


def main():
    print("Enqueing with default queue choice 'a' -> public")
    a.delay(1, 2)

    print("Enqueing with default queue choice 'b' -> private")
    b.delay(2, 3)

    print("Enqueing with overriden queue choice 'a' -> private")
    a.apply_async(args=[3, 4], queue="private")

    print("Enqueing with overriden queue choice 'b' -> public")
    b.apply_async(args=[4, 5], queue="public")


if __name__ == "__main__":
    main()
