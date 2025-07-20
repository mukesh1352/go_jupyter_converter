import torch
import torch.nn as nn
import torch.optim as optim
import matplotlib.pyplot as plt

# Dummy training data: y = 2x + 1
x_train = torch.tensor([[1.0], [2.0], [3.0], [4.0]])
y_train = torch.tensor([[3.0], [5.0], [7.0], [9.0]])

# Define the model
model = nn.Linear(in_features=1, out_features=1)

# Loss and optimizer
criterion = nn.MSELoss()
optimizer = optim.SGD(model.parameters(), lr=0.01)

# Train the model
epochs = 100
losses = []

for epoch in range(epochs):
    model.train()
    
    # Forward pass
    outputs = model(x_train)
    loss = criterion(outputs, y_train)
    losses.append(loss.item())
    
    # Backward pass
    optimizer.zero_grad()
    loss.backward()
    optimizer.step()
    
    # Logging
    if (epoch + 1) % 10 == 0:
        print(f'Epoch [{epoch+1}/{epochs}], Loss: {loss.item():.4f}')

# Plot loss
plt.plot(range(1, epochs + 1), losses)
plt.xlabel('Epoch')
plt.ylabel('Loss')
plt.title('Training Loss over Epochs')
plt.grid(True)
plt.show()

# Test prediction
model.eval()
with torch.no_grad():
    test_input = torch.tensor([[5.0]])
    prediction = model(test_input)
    print(f'\nPrediction for x=5.0: {prediction.item():.4f}')
