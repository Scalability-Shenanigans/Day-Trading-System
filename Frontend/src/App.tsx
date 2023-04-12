import { useEffect, useState } from "react";
import "./App.css";
import styled from "styled-components";
import { getAllTransactionsByUser, getBalance } from "./requests/requests";
import TransactionsList, {
  TransactionsListItemProps,
} from "./components/TransactionsList";
import LoginAlert from "./components/LoginAlert";
import Dashboard from "./components/Dashboard";

const AppContainer = styled.div`
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  box-sizing: border-box;
  max-width: 760px;
  margin: 0 auto;
  text-align: center;
  font-family: "Lato", sans-serif;
  background-color: #19232e;
  padding: 15px;
  border-radius: 10px;
`;

function App() {
  // initialize state
  const [user, setUser] = useState("");
  const [funds, setFunds] = useState(0);
  const [transactionCommitted, setTransactionCommitted] = useState(false);
  const [transactions, setTransactions] = useState<TransactionsListItemProps[]>(
    []
  );

  // fetch functions
  const fetchBalance = async () => {
    const getBalanceResponse = await getBalance({ user: user });
    const balance = parseFloat(getBalanceResponse?.data["balance"].toFixed(2));
    setFunds(balance);
  };

  const fetchUserTransactions = async () => {
    const userTransactions = await getAllTransactionsByUser({
      user: user,
    });

    setTransactions(userTransactions as TransactionsListItemProps[]);
  };

  // use Effects
  useEffect(() => {
    if (transactionCommitted === true) {
      setTransactionCommitted(false);
    }

    fetchBalance();
    fetchUserTransactions();
  }, [transactionCommitted]);

  const handleLoginSubmit = (username: string, password: string) => {
    console.log("Username:", username);
    console.log("Password:", password);
    setUser(username);
  };

  return user !== "" ? (
    <AppContainer>
      <Dashboard
        user={user}
        funds={funds}
        setFunds={setFunds}
        fetchBalance={fetchBalance}
        setTransactionCommitted={setTransactionCommitted}
        fetchUserTransactions={fetchUserTransactions}
      />
      <TransactionsList
        transactions={transactions}
        user={user}
        setCommit={setTransactionCommitted}
      />
    </AppContainer>
  ) : (
    <LoginAlert onSubmit={handleLoginSubmit} />
  );
}

export default App;
