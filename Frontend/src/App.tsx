import React, { useEffect, useState } from "react";
import "./App.css";
import TransactionForm from "./components/TransactionForm";
import styled from "styled-components";
import { addFunds, buyStock, sellStock } from "./requests/requests";
import TransactionsList, {
  TransactionsListItemProps,
} from "./components/TransactionsList";
import { AxiosResponse } from "axios";

const AddFundsText = styled.h3({
  margin: 0,
});

const AddFundsInput = styled.input({
  marginBottom: 20,
  marginRight: 5,
});

const AppContainer = styled.div`
  box-sizing: border-box;
  max-width: 760px;
  margin: auto;
  text-align: center;
  font-family: "Lato", sans-serif;
  background-color: gray;
  padding: 15px;
  margin-top: 20px;
  border-radius: 10px;
`;

const FormsContainer = styled.div`
  display: flex;
  justify-content: space-between;
`;

const TopContainer = styled.div`
  color: white;
`;

function getCurrentDateFormatted(): string {
  const now = new Date();
  const year = now.getFullYear();
  const month = (now.getMonth() + 1).toString().padStart(2, "0");
  const day = now.getDate().toString().padStart(2, "0");

  return `${year}-${month}-${day}`;
}

function App() {
  const user = "USER";

  const handleQuoteSubmit = async (stock: string) => {
    console.log("Quote submitted", stock);
  };

  const handleBuySubmit = async (stock: string, amount: number) => {
    console.log("Buy submitted", stock, amount);

    const buyRequestStatus = await buyStock({
      user: user,
      stock: stock,
      amount: amount,
    });

    if (buyRequestStatus === 200) {
      // console.log("Buy successful");

      const newTransaction = {
        type: "Buy",
        date: getCurrentDateFormatted(),
        asset: stock,
        amount: amount,
      } as TransactionsListItemProps;

      setTransactions((prevTransactions) => [
        ...prevTransactions,
        newTransaction,
      ]);
    }
  };

  const handleSellSubmit = async (stock: string, amount: number) => {
    console.log("Sell submitted", stock, amount);

    const sellRequest: AxiosResponse<any, any> | null = await sellStock({
      user: user,
      stock: stock,
      amount: amount,
    });

    console.log(sellRequest);

    if (sellRequest?.data["status"] === "success") {
      const newTransaction = {
        type: "Sell",
        date: getCurrentDateFormatted(),
        asset: stock,
        amount: amount,
      } as TransactionsListItemProps;

      setTransactions((prevTransactions) => [
        ...prevTransactions,
        newTransaction,
      ]);
    } else {
      const failureMsg = sellRequest?.data["message"];
      alert(failureMsg);
    }
  };

  const [selectedFunds, setSelectedFunds] = useState(0);
  const [funds, setFunds] = useState(0);

  const [transactions, setTransactions] = useState<TransactionsListItemProps[]>(
    []
  );

  return (
    <AppContainer>
      <TopContainer>
        <h1>Daytrading Frontend</h1>
        <h1>Logged in as {user}</h1>

        <h2>Current funds: ${funds}</h2>
        <AddFundsText>Add Funds</AddFundsText>
        <AddFundsInput
          type="number"
          value={selectedFunds}
          onChange={(e) => setSelectedFunds(Number(e.target.value))}
        />
        <button
          onClick={async () => {
            const newBalance = (await addFunds({
              user: user,
              amount: selectedFunds,
            })) as unknown as number;
            console.log("newBalance is", newBalance);
            newBalance && setFunds(newBalance);
          }}
        >
          Add funds
        </button>
      </TopContainer>

      <FormsContainer>
        <TransactionForm
          title="Quote"
          buttonText="Get Quote"
          onSubmit={(stock) => handleQuoteSubmit(stock)}
          showAmount={false}
        />
        <TransactionForm
          title="Buy"
          buttonText="Buy"
          onSubmit={(stock, amount) => handleBuySubmit(stock, amount)}
          showAmount={true}
        />
        <TransactionForm
          title="Sell"
          buttonText="Sell"
          onSubmit={(stock, amount) => handleSellSubmit(stock, amount)}
          showAmount={true}
        />
      </FormsContainer>
      <TransactionsList transactions={transactions} user={user} />
    </AppContainer>
  );
}

export default App;
