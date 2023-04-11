import React, { useEffect, useState } from "react";
import "./App.css";
import TransactionForm from "./components/TransactionForm";
import styled from "styled-components";
import {
  addFunds,
  buyStock,
  commitBuy,
  commitSell,
  getAllTransactionsByUser,
  getBalance,
  sellStock,
} from "./requests/requests";
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

const CommitButtonsContainer = styled.div`
  display: flex;
  flex-direction: horizontal;
  justify-content: center;
`;

const CommitButton = styled.button`
  margin: 0px 5px;
`;

function App() {
  const user = "USER";

  const fetchUserTransactions = async () => {
    const userTransactions = await getAllTransactionsByUser({
      user: user,
    });

    console.log("the user transactions response is", userTransactions);
    setTransactions(userTransactions as TransactionsListItemProps[]);
  };

  const fetchBalance = async () => {
    const getBalanceResponse = await getBalance({ user: user });
    const balance = parseFloat(getBalanceResponse?.data["balance"].toFixed(2));
    setFunds(balance);
  };

  const handleQuoteSubmit = async (stock: string) => {
    console.log("Quote submitted", stock);
  };

  const handleBuySubmit = async (stock: string, amount: number) => {
    console.log("Buy submitted", stock, amount);

    await buyStock({
      user: user,
      stock: stock,
      amount: amount,
    });

    fetchUserTransactions();
  };

  const handleSellSubmit = async (stock: string, amount: number) => {
    console.log("Sell submitted", stock, amount);

    await sellStock({
      user: user,
      stock: stock,
      amount: amount,
    });

    fetchUserTransactions();
  };

  const [selectedFunds, setSelectedFunds] = useState(0);
  const [funds, setFunds] = useState(0);

  const [transactionCommitted, setTransactionCommitted] = useState(false);

  useEffect(() => {
    console.log("transaction committed");

    if (transactionCommitted === true) {
      setTransactionCommitted(false);
    }

    fetchBalance();
    fetchUserTransactions();
  }, [transactionCommitted]);

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
      <CommitButtonsContainer>
        <CommitButton
          onClick={async () => {
            await commitBuy({
              user: user,
            });
            setTransactionCommitted(true);
          }}
        >
          Commit Buy
        </CommitButton>
        <CommitButton
          onClick={async () => {
            await commitSell({
              user: user,
            });
            setTransactionCommitted(true);
          }}
        >
          Commit Sell
        </CommitButton>
      </CommitButtonsContainer>
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
      <TransactionsList
        transactions={transactions}
        user={user}
        setCommit={setTransactionCommitted}
      />
    </AppContainer>
  );
}

export default App;
