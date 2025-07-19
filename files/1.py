import pandas as pd
import numpy as np

# Create a simple DataFrame
data = {
    'Name': ['Alice', 'Bob', 'Charlie', 'David'],
    'Math Score': [85, 90, 78, 92],
    'Science Score': [88, 76, 93, 85]
}

df = pd.DataFrame(data)

# Add a new column for average score using NumPy
df['Average Score'] = np.mean(df[['Math Score', 'Science Score']], axis=1)

# Display the DataFrame
print("Student Scores:\n")
print(df)

# Find the student with the highest average score
top_student = df.loc[df['Average Score'].idxmax()]
print("\nTop Student:")
print(top_student)
