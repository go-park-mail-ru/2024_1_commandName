# вот часть нашего кода на Python


import json
import logging
import random
import string

from locust import HttpUser, task, between, TaskSet


class WebsiteTasks(TaskSet):
    def on_start(self):
        pass

    def generate_username(self, length=10):
        letters = string.ascii_letters
        return ''.join(random.choice(letters) for i in range(length))

    def create_private_chat(self, user_id):
        print(user_id)
        response = self.client.post("/createPrivateChat", json={
            "user_id": user_id
        })

        if response.status_code != 200:
            logging.error(f"Failed to create chat for user {user_id}: {response.status_code} {response.text}")
            try:
                error_response = response.json()
                logging.error(f"Error response: {json.dumps(error_response, indent=2)}")
            except json.JSONDecodeError:
                logging.error("Response is not JSON")
                logging.error(response.text)
        else:
            logging.info(f"Successfully created private chat for user {user_id}")

    def logout_user(self):
        response = self.client.post("/logout")
        if response.status_code != 200:
            logging.error(f"Failed to logout: {response.status_code} {response.text}")
            try:
                error_response = response.json()
                logging.error(f"Error response: {json.dumps(error_response, indent=2)}")
            except json.JSONDecodeError:
                logging.error("Response is not JSON")
                logging.error(response.text)
        else:
            logging.info("Successfully logged out")

    @task
    def create_chats_task(self):
        counter = 2000
        userArr = ["NJzsNBStPp", "DbAzZXGvWs", "RImteragvh", "cweHGubJyu", "JpodqLEZSC", "TgRRpKmmue", "RcWYiWIKbn", "FiqZAROFjP", "ArtemKa", "iPvSrJwLCn"]
        for user in userArr:
            self.login_user(user)
            for i in range(counter, counter+10000):
                self.create_private_chat(user_id=i)
            self.logout_user()
            counter += 10000
            print("counter", counter)


    def login_user(self, username):
        password = "Demouser123!"

        headers = {
            "Content-Type": "application/json"
        }

        response = self.client.post("/login", json={
            "username": username,
            "password": password
        }, headers=headers)

        if response.status_code != 200:
            logging.error(f"Failed to register: {response.status_code} {response.text}")
            try:
                error_response = response.json()
                logging.error(f"Error response: {json.dumps(error_response, indent=2)}")
            except json.JSONDecodeError:
                logging.error("Response is not JSON")
                logging.error(response.text)
            return False
        else:
            logging.info(f"Successfully registered: {username}")
            return True




class WebsiteUser(HttpUser):
    tasks = [WebsiteTasks]
    wait_time = between(5, 15)


if __name__ == "__main__":
    import locust

    locust.run_single_user(WebsiteUser)
