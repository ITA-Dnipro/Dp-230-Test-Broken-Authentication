# Dp-230-Test-Broken-Authentication

### Brute force scanner v0.2

**Example**

Input: Network tab, and request payload
![1](https://user-images.githubusercontent.com/19240229/176669437-1c3eebd5-e206-468e-8615-a6255c3c4139.png)


Output: log with mixed corr and wrong answers
![2](https://user-images.githubusercontent.com/19240229/176669465-3a4060df-144a-4dcd-9090-731577367ed6.png)



**Core logic:** 

- **Input**: data(endpoint) provided manually
- Dictionaries generated in txt file(try.txt will be used for password attempts)
- _TODO_: Endpoint read from .env file
- Logic for POST request and logic for bruteforce the target
- **Output**: "if pass correct" write to correct_password.txt.(Optional show log in console)
- _TODO_: test on local env
- _TODO_: test using online resource

**Features:**

- Provides logging to separate file and console
- Uses a combination of password that fit the password rule(minimum length of 8, has to include uppercase and number)
- Checks whether the status response success or error
- Tracking system the brute force can be paused and continued later(used plain txt)

**Pros:**

- handled unified makefile implementation
- used best practice in lintering and pre commit configuration checks

**Cons:**

- Hardcoded password generator
- Test coverage 0
- No automation logic implemented for further crawler usage
- Outpout isn't unified for proceeding to report service
