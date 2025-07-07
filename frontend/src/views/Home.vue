<template>
  <div class="home">
    <div class="hero">
      <div class="hero-content">
        <h1 class="hero-title">
          Rejoignez des <span class="gradient-text">communautés</span> 
          passionnantes
        </h1>
        <p class="hero-subtitle">
          Découvrez, créez et participez à des groupes qui partagent vos centres d'intérêt.
          Connectez-vous avec des personnes qui vous ressemblent.
        </p>
        <div class="hero-actions">
          <router-link to="/groups" class="btn btn-primary">
            Découvrir les groupes
          </router-link>
          <router-link to="/groups/create" class="btn btn-secondary">
            Créer un groupe
          </router-link>
        </div>
      </div>
      <div class="hero-image">
        <img src="https://images.pexels.com/photos/3184292/pexels-photo-3184292.jpeg?auto=compress&cs=tinysrgb&w=600" alt="Community" />
      </div>
    </div>
    
    <div class="features">
      <div class="container">
        <h2 class="section-title">Pourquoi nous rejoindre ?</h2>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon">🌐</div>
            <h3>Communautés diverses</h3>
            <p>Trouvez des groupes sur tous les sujets : technologie, art, sport, cuisine et bien plus.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">📅</div>
            <h3>Événements exclusifs</h3>
            <p>Participez à des événements organisés par les membres de votre communauté.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">💬</div>
            <h3>Discussions engageantes</h3>
            <p>Partagez vos idées, posez des questions et apprenez des autres membres.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">🔒</div>
            <h3>Groupes privés</h3>
            <p>Créez des espaces intimes pour vos proches ou votre équipe de travail.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// Home page component
</script>

<style scoped>
.home {
  min-height: 100vh;
}

.hero {
  display: flex;
  align-items: center;
  min-height: 80vh;
  padding: 0 20px;
  max-width: 1200px;
  margin: 0 auto;
  gap: 60px;
}

.hero-content {
  flex: 1;
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 800;
  line-height: 1.1;
  margin-bottom: 20px;
  color: #fff;
}

.gradient-text {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 1.25rem;
  color: #ccc;
  line-height: 1.6;
  margin-bottom: 30px;
}

.hero-actions {
  display: flex;
  gap: 16px;
}

.hero-image {
  flex: 1;
  max-width: 500px;
}

.hero-image img {
  width: 100%;
  height: auto;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.features {
  padding: 80px 20px;
  background: #111;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.section-title {
  font-size: 2.5rem;
  font-weight: 700;
  text-align: center;
  margin-bottom: 50px;
  color: #fff;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 30px;
}

.feature-card {
  background: #1a1a1a;
  padding: 30px;
  border-radius: 12px;
  text-align: center;
  border: 1px solid #333;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
}

.feature-icon {
  font-size: 3rem;
  margin-bottom: 20px;
}

.feature-card h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 15px;
  color: #fff;
}

.feature-card p {
  color: #ccc;
  line-height: 1.6;
}

.btn {
  padding: 16px 32px;
  border-radius: 12px;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
}

.btn-primary {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
}

.btn-primary:hover {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
  transform: translateY(-2px);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
}

@media (max-width: 768px) {
  .hero {
    flex-direction: column;
    text-align: center;
    gap: 40px;
    padding: 40px 20px;
  }
  
  .hero-title {
    font-size: 2.5rem;
  }
  
  .hero-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .hero-actions .btn {
    width: 100%;
    max-width: 300px;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style>