### Installation

1. Clone the repository:

    git clone https://github.com/ar-vi-nd/gotodo.git

2. Install dependencies:

    go mod tidy


3. Create a `.env` file in the root directory of the project and add your MongoDB URI and server port:

    MONGODB_URI=mongodb://localhost:27017/todoapp
    PORT=3000


4. Run the application:

    go run main.go


5. Your server should be running at `http://localhost:3000`.



### API Endpoints

#### List All Todo-> `/todos`   GET

#### Create a Todo-> `/todos`   POST  
 **Body**: JSON
    ```json
    {
        "title": "Example Title",
        "description": "Example Description",
        "completed": false
    }

#### Update a Todo by ID
- **URL**: `/todos/{id}`
- **Method**: `PUT`
- **Body**: JSON
    ```json
    {
        "title": "Updated Title",
        "description": "Updated Description",
        "completed": true
    }
    ```

#### Delete a Todo by ID
- **URL**: `/todos/{id}`
- **Method**: `DELETE`
