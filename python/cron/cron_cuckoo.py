import cron
import datetime


for _ in cron.cron(10, 5):
    print(datetime.datetime.now().strftime('%Y/%m/%d %H:%M:%S'), 'main: cuckoo')
