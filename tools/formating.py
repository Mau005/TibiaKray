book = '''
Upload image
Characters
Messages
Favorite Images
Friends
Notifications
My account
Setting
My profile
Language
Change Password
New Password
Repeat New Password
Save Changes
'''

newBook  = book.split("\n")

procesing = []

for elements in newBook:
    if elements == "":
        continue
    procesing.append(elements.replace(" ", "").lower())


newListCopy = ""

for elements in procesing:
    newListCopy += elements+"\n"

print(newListCopy)

