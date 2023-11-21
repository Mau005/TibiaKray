book = '''
Tell us your experience
Upload Image
Authorize Post
Delete post
To comment you need to be registered
Published
Written by
View
Comments
Voted
Post your comment
Post Comment
Message
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

