# Digital Closet Web Application

## Functional Requirements (MVP)
- **User Account Management**
  - Basic user registration and login.
  - Password reset.
  - Minimal user profile management (username, password, email).
- **Clothing Management**
  - Ability for users to upload photos of clothing via file upload.
  - Users can create and categorize their clothing items using custom categories.
  - Functions to edit and delete clothing items.
- **Outfit Assembly**
  - A sandbox interface where users can drag and drop clothing items.
  - The interface will save the position of items on the space.
- **Image Processing**
  - Automatic cropping to isolate clothing from the background in uploaded photos.
- **Search and Filter**
  - Initially, no search or filter functionality.
- **Social Features**
  - No social features in the MVP.
- **Notifications**
  - No notifications in the MVP.

## Non-Functional Requirements (MVP)
- **Performance**
  - Fast response times, especially critical for image uploads and interactions in the sandbox.
- **Scalability**
  - Ability to handle growth in user numbers and data volume efficiently.
- **Security**
  - Basic security measures considering the nature of the data involved (clothing images, email addresses).
- **Usability**
  - Intuitive and user-friendly interface accessible from various devices.
- **Reliability and Maintainability**
  - High application availability with minimal downtime.
  - Well-documented code for easy maintenance and future updates.

## Architecture Design
- **Frontend**
  - **Framework/Technology:** React
  - **Design Tools:** Figma for UI/UX design.
- **Backend**
  - **Framework/Technology:** Flask (Python)
  - **Database:** PostgreSQL
  - **File Storage:** AWS S3
- **Infrastructure Setup**
  - **Hosting:** AWS EC2 or Elastic Beanstalk
  - **Development Environment:** Docker for containerized development.

## Tech Stack
- **Frontend:** React
- **Backend:** Flask (Python)
- **Database:** PostgreSQL
- **File Storage:** AWS S3
- **Hosting:** AWS services (EC2 or Elastic Beanstalk)
- **Development Environment:** Docker

## API Design
- REST API will be documented using Swagger to ensure clarity and ease of use.

## Testing Frameworks
- **Backend Testing:** pytest for Python/Flask.
- **Frontend Testing:** Jest combined with React Testing Library.
- **End-to-End Testing:** Selenium or Cypress.

## Deployment Strategy
- Use AWS Elastic Beanstalk for simplified application deployment and management.
- Implement CI/CD pipelines using GitHub Actions for automated testing and deployment processes.

## Project Management
- Git for version control.
- Mainly an individual project or small team collaboration without extensive project management tools.

## Database Schema (Initial Sketch)
- **Users:** Store user details such as username, password, and email.
- **Clothing Items:** Attributes include category, image URL, and user association.
- **Categories:** Custom categories created by users.
