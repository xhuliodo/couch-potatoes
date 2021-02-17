from PIL import Image
import os

directory= r'C:\Users\xhulio.doda\Desktop\cp\dataset\06-download-poster-pics\img'

for filename in os.listdir(directory):
    filepath = '../06-download-poster-pics/img/'+filename
    image = Image.open(filepath)
    image.thumbnail((450, 650))
    image.save('thumbnail_'+filename)
    image.save('thumbnail_'+filename)
    size = os.path.getsize(filepath)