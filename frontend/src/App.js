import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';
import ProductList from './components/ProductList';
import Cart from './components/Cart';
import CreateProductForm from './components/CreateProductForm';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8082';

const App = () => {
  const [products, setProducts] = useState([]);
  const [cart, setCart] = useState([]);

  // Fetch products from the backend
  useEffect(() => {
    axios.get(API_URL + '/product')
      .then((response) => {
          setProducts(response.data.Products);
      })
      .catch((error) => console.error('Error fetching products:', error));
  }, []);

  // Add a product to the cart
  const addToCart = (product) => {
    setCart((prevCart) => [...prevCart, product]);
  };

  // Remove a product from the cart
  const removeFromCart = (productId) => {
    setCart((prevCart) => prevCart.filter((item) => item.id !== productId));
  };

  return (
    <div className="App">
      <h1>Online Shop</h1>
      <ProductList products={products} addToCart={addToCart} />
      <Cart cart={cart} removeFromCart={removeFromCart} />
      <CreateProductForm setProducts={setProducts} />
    </div>
  );
};

export default App;
