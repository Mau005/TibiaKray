from datetime import datetime, timedelta
price = 9000

month_list_name =["", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"]
ultima_fecha =	"5 de marzo"
date_init = datetime(2023, 3, 5) #Fecha en que se saco todo

book =	'''
Febrero: [16, 19, 23, 26]
Marzo: [2, 5]
Total Adeudado: 460.000
'''

book_abonos = '''
Transferencia 12 de marzo 2023: 50.000
Transferencia 22 de marzo 2023: 50.000
Transferencia 18 de junio 2023: 50.000
Transferencia 2 de julio 2023: 50.000
Efectivo 23 de julio 2023: 50.000
Transferencia  5 de septiembre 2023: 200.000
Transferencia 1 de octubre 2023: 50.000 
'''
total_abonos = (50000 *	6) + 200000
month_list  =	{}

def add_secuency(key) -> list:
	if month_list .get(key) is None:
		month_list.update({key:[]})
		
	return month_list.get(key)

days_target =	4
secuency = True
while True:
	date_init += timedelta(days=days_target) #domingos
	add_secuency(date_init.month)
	month_list[date_init.month].append(date_init.day)
	
	if secuency:
		days_target = 3
		secuency =False
	else:
		days_target =	4
		secuency=True
		
	if date_init.month == 11 and date_init.day >= 5:
		break

print(book)
print(book_abonos)
total = 460000
for elements in month_list:
	content =	month_list[elements]
	print(f"{month_list_name[elements]}: {content}")
	for values in content:
		total += price
dias_faltados =4

print(f"Dias Faltados: {dias_faltados}")
print(f"Total Abonos: {total_abonos}")
print(f"Total: {total}")
print(f"Deuda: {total - total_abonos - (dias_faltados * price)}")
