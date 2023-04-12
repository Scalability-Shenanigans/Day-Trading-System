import React, { useState } from "react";
import styled from "styled-components";

interface TransactionFormProps {
  title: string;
  buttonText: string;
  onSubmit: (stock: string, amount: number) => void;
  showAmount: boolean;
  quotePrice?: number | null;
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
  box-sizing: border-box;
`;

const FieldContainer = styled.div`
  display: flex;
  flex-direction: column;
  margin-bottom: 10px;
  box-sizing: border-box;
`;

const Label = styled.label`
  margin-bottom: 5px;
`;

const Input = styled.input`
  // margin-bottom: 15px;
`;

const Quote = styled.h3`
  font-weight: bold;
  text-align: center;
`;

const Button = styled.button`
  outline: none;
  color: white;
  background: transparent;
  padding: 5px 15px;
  border: 2px solid white;
  border-radius: 5px;
  font-weight: bold;

  &:hover {
    color: white;
    background: #535354;
    border-color: white;
  }
`;

const TransactionForm: React.FC<TransactionFormProps> = ({
  title,
  buttonText,
  onSubmit,
  showAmount,
  quotePrice,
}) => {
  const [stock, setStock] = useState("");
  const [amount, setAmount] = useState<number | null>(null);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSubmit(stock, amount ?? 0);
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
              value={amount ?? ""}
              onChange={(e) => setAmount(Number(e.target.value))}
            />
          </FieldContainer>
        )}
        {quotePrice && (
          <FieldContainer>
            <Quote>${quotePrice}</Quote>
          </FieldContainer>
        )}
      </div>

      <Button type="submit">{buttonText}</Button>
    </FormContainer>
  );
};

export default TransactionForm;
