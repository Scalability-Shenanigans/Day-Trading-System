import React, { useState } from "react";
import styled from "styled-components";

interface LoginAlertProps {
  onSubmit: (username: string, password: string) => void;
}

const Wrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: calc(100vh - 120px);
  background-color: transparent;
`;

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin: auto;
`;

const LoginForm = styled.form`
  background-color: #34495e;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 300px;
`;

const InputField = styled.div`
  display: flex;
  flex-direction: column;
  margin-bottom: 1rem;
`;

const Label = styled.label`
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 5px;
  color: #ecf0f1;
`;

const Input = styled.input`
  border: 1px solid #7f8c8d;
  border-radius: 4px;
  font-size: 14px;
  padding: 8px;
  outline: none;
  background-color: #2c3e50;
  color: #ecf0f1;
  &:focus {
    border-color: #4a90e2;
  }
`;

const Button = styled.button`
  background-color: #4a90e2;
  color: white;
  font-weight: bold;
  font-size: 14px;
  border: none;
  border-radius: 4px;
  padding: 8px;
  cursor: pointer;
  &:hover {
    background-color: #2080e8;
  }
`;

const LoginAlert: React.FC<LoginAlertProps> = ({ onSubmit }) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = (e: any) => {
    e.preventDefault();
    onSubmit(username, password);
  };

  return (
    <Wrapper>
      <Container>
        <LoginForm onSubmit={handleSubmit}>
          <InputField>
            <Label htmlFor="username">Username</Label>
            <Input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </InputField>
          <InputField>
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </InputField>
          <Button type="submit">Log in</Button>
        </LoginForm>
      </Container>
    </Wrapper>
  );
};

export default LoginAlert;
