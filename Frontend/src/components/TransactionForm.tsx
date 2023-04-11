import React, { useState } from "react";
import styled from "styled-components";

interface TransactionFormProps {
  title: string;
  buttonText: string;
  onSubmit: (stock: string, amount: number) => void;
  showAmount: boolean;
}

const FormContainer = styled.form`
  padding: 15px 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  margin: 20px 0px;
  text-align: left;

  background-color: #353638;
  border-radius: 5px;
  align-self: stretch;
  color: white;
`;

const FieldContainer = styled.div`
  display: flex;
  flex-direction: column;
  margin-bottom: 10px;
`;

const Label = styled.label`
  margin-bottom: 5px;
`;

const Input = styled.input`
  // margin-bottom: 15px;
`;

const Button = styled.button({
  // border: "none",
  outline: "none",
  color: "white",
  background: "transparent",
  padding: "5px 15px",
  border: "2px solid white",
  borderRadius: 5,
  fontWeight: "bold",
});

const TransactionForm: React.FC<TransactionFormProps> = ({
  title,
  buttonText,
  onSubmit,
  showAmount,
}) => {
  const [stock, setStock] = useState("");
  const [amount, setAmount] = useState(0);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSubmit(stock, amount);
  };

  return (
    <FormContainer onSubmit={handleSubmit}>
      <div>
        <h3>{title}</h3>
        <FieldContainer>
          <Label htmlFor="stock">Stock: </Label>
          <Input
            type="text"
            name="stock"
            value={stock}
            onChange={(e) => setStock(e.target.value)}
          />
        </FieldContainer>
        {showAmount && (
          <FieldContainer>
            <Label htmlFor="amount">Amount: </Label>
            <Input
              type="number"
              name="amount"
              value={amount}
              onChange={(e) => setAmount(Number(e.target.value))}
            />
          </FieldContainer>
        )}
      </div>

      <Button type="submit">{buttonText}</Button>
    </FormContainer>
  );
};

export default TransactionForm;