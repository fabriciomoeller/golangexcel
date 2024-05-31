from pathlib import Path
import pandas as pd

df = pd.read_excel(
  io = Path.cwd() / 'volumetria.xlsx',
  engine='openpyxl',
  sheet_name='LGM',
  skiprows=1,
  usecols='B:D',
  nrows=401
)
df.head()
