Here's a clear and organized breakdown of your "Profile" feature specification into sections and bullet points for easier implementation and understanding:

---

### 🔹 Profile Overview

Every user profile must include:

1. **User Information**
   Display all information provided during registration, *except* the password:

   * Email
   * First Name
   * Last Name
   * Date of Birth
   * Nickname
   * About Me
   * Avatar URL
   * Public/Private status

2. **User Activity**

   * List of **all posts** made by the user.

3. **Social Connections**

   * **Followers**: Users who follow this profile.
   * **Following**: Users this profile is following.

---

### 🔹 Profile Visibility Types

1. **Public Profile**

   * Anyone on the platform can view the user's information, posts, followers, and following list.

2. **Private Profile**

   * Only **followers** of the user can see:

     * User info
     * Posts
     * Followers / Following lists

---

### 🔹 Profile Owner Controls

When a user views *their own* profile:

* An option must be available to:

  * **Toggle profile visibility** (Public ↔ Private)

---

Let me know if you'd like this broken down further by database schema, API route design, or frontend rendering logic.
