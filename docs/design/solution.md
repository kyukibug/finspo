# Solution Documentation

## Author
- **Name**: Alex N.
- **Contact Information**: axni@umich.edu

## Version History
| Version | Date       | Description                                    |
|---------|------------|------------------------------------------------|
| 1.0     | 2024-07-08 | Initial version of the solution documentation. |

## Last Modified
- **Date**: 2024-07-08

## Purpose
This document provides a comprehensive overview of the digital closet application. It details the system architecture, requirements, detailed design, testing, and deployment strategies. The purpose is to serve as a guide for development and maintenance and ensure alignment with business goals and user needs.

## Architecture Design
![Architecture Design](tbd)

## Requirements
### Functional
- **User Registration and Authentication**
  - Users can register using their Google account to authenticate.
  - Secure authentication mechanism using OAuth.
  - Users can edit their username.

- **Clothing Management**
  - Users can upload images of their clothing items.
  - Automatic image cropping to extract clothing from uploaded images. 
  - Users can create custom categories to categorize clothing items.
  - Users can edit and delete existing clothing items.
  - Users can filter by category
 

- **Sandbox Interface**
  - Users can drag and drop clothing items to arrange them in a virtual space.
  - The position of each item within the sandbox can be saved and retrieved.



### Non-Functional
- Performance: Application should handle multiple user interactions swiftly.
- Scalability: Must support scaling to handle growth in user base and data.
- Security: Basic security practices must be implemented, considering the sensitivity of user data.

--- 
#### Out of Scope
- Social sharing features.
- Advanced image editing tools.
- Real-time collaboration on outfit designs.

## Detailed Design
### Endpoints
- List of primary RESTful endpoints:
  - `POST /users`: Register a new user.
  - `GET /clothes`: Retrieve user's clothing items.
  - `PUT /clothes/{id}`: Update specific clothing item details.

### Database Schema
The following bullets outline the structure of the database:

- **Users**
  - `id`: Primary key
  - `username`: User's chosen name
  - `email`: User's email address
  - `google_id`: Identifier provided by Google OAuth

- **Clothing Items**
  - `id`: Primary key
  - `user_id`: Foreign key linked to Users
  - `category_id`: Foreign key linked to Categories
  - `image_url`: URL path to the stored image

- **Categories**
  - `id`: Primary key
  - `user_id`: Foreign key linked to Users
  - `name`: Name of the category

- **Sandbox**
  - `id`: Primary key
  - `user_id`: Foreign key linked to Users

- **Sandbox Positions**
  - `id`: Primary key
  - `sandbox_id`: Foreign key linked to Sandbox
  - `clothing_item_id`: Foreign key linked to Clothing Items
  - `position_x`: X coordinate of the item in the sandbox
  - `position_y`: Y coordinate of the item in the sandbox


## Testing Strategy
- **Unit Testing**: Focus on backend logic using pytest.
- **Integration Testing**: Ensure all components work together as expected.
- **UI Testing**: Automated testing of the frontend with tools like Selenium.

## Deployment Strategy
- **Initial Setup**: Use Docker containers to manage application environments.
- **CI/CD Pipeline**: Implement using GitHub Actions for automated testing and deployment to AWS.

## Additional Notes
- Any further details or considerations can be added here.

