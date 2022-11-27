from celery import Celery

app = Celery("tasks", broker="redis://redis:6379/0")


@app.task(queue="public")
def a(x, y):
    print(f"Printing from fn 'a': {x + y}")


@app.task(queue="private")
def b(x, y):
    print(f"Printing from fn 'b': {x * y}")
