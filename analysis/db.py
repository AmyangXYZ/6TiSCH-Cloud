import mysql.connector

db = mysql.connector.connect(host="localhost",
                             user="root",
                             passwd="0311",
                             database="6tisch")

def sqlFetch(sqlStr):
    print(sqlStr)
    cursor = db.cursor(dictionary=True)
    cursor.execute(sqlStr)
    return cursor.fetchall()

def sqlInsert(sqlstr):
    print(sqlstr)
    cursor = db.cursor()
    cursor.execute(sqlstr)

def createTable():
    cursor = db.cursor()
    cursor.execute('''CREATE TABLE IF NOT EXISTS NW_DATA_SET_NOISE_LEVEL(
        TIMESTAMP BIGINT NOT NULL,
        GATEWAY_NAME VARCHAR(64) NOT NULL,
        SENSOR_ID SMALLINT NOT NULL,
        RES_CATEGORY VARCHAR(64),
        NOISE_LEVEL SMALLINT,
        GPS_LAT DOUBLE,
        GPS_LONG DOUBLE);''')
    cursor.execute('''CREATE TABLE IF NOT EXISTS NW_DATA_SET_ANALYSIS (
        DATETIME DATETIME NOT NULL,
        TIMESTAMP BIGINT NOT NULL,
        GATEWAY_NAME VARCHAR(64),
        SENSOR_ID SMALLINT NOT NULL,
        RESULT TEXT); ''')


# createTable()