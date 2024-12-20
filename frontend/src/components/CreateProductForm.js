import React, { useState } from 'react';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8082';

const CreateProductForm = ({ setProducts }) => {
  const [newProduct, setNewProduct] = useState({
    name: '',
    description: '',
    image_url: '',
    price: '',
    weight: 1,
  });

  const createProduct = () => {
    axios.post(API_URL + '/product', newProduct)
      .then((response) => {
        setProducts((prevProducts) => [...prevProducts, response.data]);
        setNewProduct({ name: '', description: '', image_url: '', price: '', weight: 1 });
      })
      .catch((error) => console.error('Error creating product:', error));
  };

  return (
    <div className="new-product-form">
      <h2>Create New Product</h2>
      <input
        type="text"
        placeholder="Name"
        value={newProduct.name}
        onChange={(e) => setNewProduct({ ...newProduct, name: e.target.value })}
      />
      <input
        type="text"
        placeholder="Description"
        value={newProduct.description}
        onChange={(e) => setNewProduct({ ...newProduct, description: e.target.value })}
      />
      <input
        type="text"
        placeholder="Image URL"
        value={newProduct.image_url}
        onChange={(e) => setNewProduct({ ...newProduct, image_url: e.target.value })}
      />
      <input
        type="number"
        placeholder="Price"
        value={newProduct.price}
        onChange={(e) => setNewProduct({ ...newProduct, price: parseFloat(e.target.value) })}
      />
      <input
        type="number"
        placeholder="Weight"
        value={newProduct.weight}
        onChange={(e) => setNewProduct({ ...newProduct, weight: parseFloat(e.target.value) })}
      />
      <button onClick={createProduct}>Create Product</button>
    </div>
  );
};

export default CreateProductForm;
